package main

import (
	"context"
	"errors"
	"fmt"
	"net"
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
	runGrpcServer(ctx, waitGroup, config, store, nil)

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

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
	store db.Store,
	taskDistributor any,
) {
	server, err := server.New(config, store, taskDistributor)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create server")
		fmt.Println("cannot create server")
	}

	// gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer()
	pb.RegisterHavlabsAuthServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create listener")
		fmt.Println("cannot create listener")
	}

	waitGroup.Go(func() error {
		// log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
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
		fmt.Println("graceful shutdown gRPC server")

		return nil
	})
}
