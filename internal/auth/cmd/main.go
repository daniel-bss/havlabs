package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/daniel-bss/havlabs/auth/pb"
	"github.com/daniel-bss/havlabs/auth/server"
	"github.com/daniel-bss/havlabs/auth/utils"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// TODO: setup DB

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)
	runGrpcServer(ctx, waitGroup, config, nil, nil)

	err = waitGroup.Wait()
	if err != nil {
		// log.Fatal().Err(err).Msg("error from wait group")
		fmt.Println("error from wait group")
	}
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config utils.Config,
	store any,
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
