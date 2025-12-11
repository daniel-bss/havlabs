package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/server"
	"github.com/daniel-bss/havlabs/internal/news/utils"
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
	config, err := utils.LoadConfig(".")
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot load config")
	}

	// TODO: if development then:
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGRPCServer(ctx, waitGroup, config)
}

func runGRPCServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
	// store db.Store,
	// taskDistributor any,
) {
	// validator, err := protovalidate.New()
	// if err != nil {
	// 	fmt.Println("failed to create protovalidate validator")
	// 	return
	// }
	// recoveryOpts := []grpc_recovery.Option{
	// 	grpc_recovery.WithRecoveryHandler(func(interface{}) error {
	// 		return utils.InternalServerError()
	// 	}),
	// }

	service, err := server.NewGRPCService(config)
	if err != nil {
		fmt.Println("cannot create gRPC service", err.Error())
		// log.Fatal().Err(err).Msg("cannot create gRPC service")
	}

	grpcServer := grpc.NewServer(
	// grpc.ChainUnaryInterceptor(
	// 	grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
	// 	// server.CustomInterceptor(),
	// 	protovalidate_middleware.UnaryServerInterceptor(validator),
	// ),
	) // runs in order

	pb.RegisterHavlabsNewsServer(grpcServer, service)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		// log.Info().M
		// .Printf("start gRPC server at %s\n", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				fmt.Println(err.Error())
				fmt.Println("KASMDKSAD")
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
