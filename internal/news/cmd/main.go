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
	"github.com/daniel-bss/havlabs/internal/news/server"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
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

	waitGroup, ctx := errgroup.WithContext(ctx)
	runGRPCServer(ctx, waitGroup, config)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGRPCServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
	// store db.Store,
	// taskDistributor any,
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

	service, err := server.NewGRPCService(config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			// server.CustomInterceptor(),
			protovalidate_middleware.UnaryServerInterceptor(validator),
		),
	)

	pb.RegisterHavlabsNewsServer(grpcServer, service)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("serving gRPC server at %s", listener.Addr().String())

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
