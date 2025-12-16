package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/server"
	"github.com/daniel-bss/havlabs/internal/auth/token"
	"github.com/daniel-bss/havlabs/internal/auth/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// TODO: if development then:
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.GetDBSource())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.GetDBSource())

	store := db.NewStore(connPool)

	waitGroup, ctx := errgroup.WithContext(ctx)
	runGRPCServer(ctx, waitGroup, config, store, nil)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("auth DB successfully migrated")
}

func runGRPCServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
	store db.Store,
	taskDistributor any,
) {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create protovalidate validator")
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(interface{}) error {
			log.Warn().Err(err).Msg("encountered Internal Server Error")
			return utils.InternalServerError()
		}),
	}

	service, err := server.NewGRPCService(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			server.CustomInterceptor(),
			protovalidate_middleware.UnaryServerInterceptor(validator),
		),
	) // runs in order

	pb.RegisterHavlabsAuthServer(grpcServer, service)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("serving auth service at %s\n", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				log.Error().Err(err).Msg("ErrServerStopped")
				return nil
			}

			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		log.Info().Msg("graceful shutdown gRPC server")
		grpcServer.GracefulStop()

		return nil
	})
}

func runHTTPServer(config utils.Config) {
	tokenMaker, _ := token.NewJWTMaker()

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		pubKey := tokenMaker.PublicKey()

		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"kid": "oneandonly-jwk-until-further-implementation", // OPTIONAL UNTIL MULTIPLE JWKs
			"kty": "RSA",
			"alg": "RS384",
			"use": "sig",
			"n":   base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pubKey.E)).Bytes()),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to encode JWKS")
		}
	})

	port := ":8081"
	log.Info().Msgf("Server starting on port %s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}
