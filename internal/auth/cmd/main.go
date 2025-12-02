package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/pb"
	"github.com/daniel-bss/havlabs/internal/auth/server"
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
	runJWKSServer(ctx, waitGroup, config)

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
	s, err := server.NewGRPC(config, store, taskDistributor)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create server")
		fmt.Println("cannot create server")
	}

	gprcLogger := grpc.UnaryInterceptor(server.GrpcHehe)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterHavlabsAuthServer(grpcServer, s)
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

func runJWKSServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
) {
	listener := ":8080"
	mux := http.NewServeMux()
	var server *http.Server

	waitGroup.Go(func() error {
		fmt.Printf("start JWKS server at %s\n", listener)

		mux.HandleFunc("/jwks.json", func(w http.ResponseWriter, r *http.Request) {
			kid := "my-key-id-1"

			jwk := map[string]any{
				"kty": "RSA",
				"kid": kid,
				"alg": "RS256",
				"use": "sig",
				// "n":   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
				// "e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
			}

			json.NewEncoder(w).Encode(map[string]any{
				// "data": []string{"iPhone", "MacBook", "iPad"},
				"keys": []any{jwk},
			})
		})

		server = &http.Server{
			Addr:    listener,
			Handler: mux,
		}

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on %s: %v\n", listener, err)
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Server shutdown failed: %v\n", err)
		}
		fmt.Println("graceful shutdown JWKS server")

		return nil
	})
}
