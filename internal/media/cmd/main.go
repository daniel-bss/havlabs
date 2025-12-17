package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/client"
	"github.com/daniel-bss/havlabs/internal/media/server"
	"github.com/daniel-bss/havlabs/internal/media/usecases"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/golang-migrate/migrate/v4"                     // important for golang-migrate
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // important for golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // important for golang-migrate
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

	// config
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	// context
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// connection pool
	connPool, err := pgxpool.New(ctx, config.GetDBSource())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	// db migration
	runDBMigration(config.MigrationURL, connPool)

	// waitgroup
	waitGroup, ctx := errgroup.WithContext(ctx)
	runGRPCServer(ctx, waitGroup, config)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runDBMigration(migrationURL string, connPool *pgxpool.Pool) {
	connString := connPool.Config().ConnConfig.ConnString()
	migration, err := migrate.New(migrationURL, connString)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("media DB successfully migrated")
}

func runGRPCServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
	// taskDistributor any,
) {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create protovalidate validator")
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(any) error {
			log.Warn().Err(err).Msg("encountered Internal Server Error")
			return utils.InternalServerError()
		}),
	}

	// init minio, store, usecase, server
	minioManager := client.NewMinioManager(config)
	usecase := usecases.New(minioManager)
	service, err := server.NewGRPCService(config, usecase)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			server.CustomInterceptor(),
			protovalidate_middleware.UnaryServerInterceptor(validator),
		),
	)

	pb.RegisterHavlabsMediaServer(grpcServer, service)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("serving media service at %s", listener.Addr().String())

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

		log.Info().Msg("graceful shutdown media gRPC server")
		grpcServer.GracefulStop()

		return nil
	})
}
