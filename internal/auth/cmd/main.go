package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/server"
	"github.com/daniel-bss/havlabs/internal/auth/token"
	"github.com/daniel-bss/havlabs/internal/auth/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
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
	config, err := utils.LoadConfig(".")
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot load config")
		fmt.Println("cannot load config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.GetDBSource())
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot connect to db")
		fmt.Println("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.GetDBSource())

	store := db.NewStore(connPool)

	waitGroup, ctx := errgroup.WithContext(ctx)
	runGRPCServer(ctx, waitGroup, config, store, nil)

	err = waitGroup.Wait()
	if err != nil {
		// log.Fatal().Err(err).Msg("error from wait group")
		fmt.Println("error from wait group")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create new migrate instance")
		fmt.Println("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		// log.Fatal().Err(err).Msg("failed to run migrate up")
		fmt.Println("failed to run migrate up")
	}

	// log.Info().Msg("db migrated successfully")
	fmt.Println("db migrated successfully")
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
		fmt.Println("failed to create protovalidate validator")
		return
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(interface{}) error {
			return utils.InternalServerError()
		}),
	}

	service, err := server.NewGRPCService(config, store, taskDistributor)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create server")
		fmt.Println("cannot create gRPC service")
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
		// log.Fatal().Err(err).Msg("cannot create listener")
		fmt.Println("cannot create listener")
	}

	waitGroup.Go(func() error {
		fmt.Printf("start gRPC server at %s\n", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			// log.Error().Err(err).Msg("gRPC server failed to serve")
			fmt.Println("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		fmt.Println("graceful shutdown gRPC server")
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
			fmt.Println("failed to encode JWKS")
		}
	})
	port := ":8081"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println(http.ListenAndServe(port, mux))
}
