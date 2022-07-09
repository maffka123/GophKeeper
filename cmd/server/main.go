package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/app"
	basecfg "github.com/maffka123/GophKeeper/internal/config"
	"github.com/maffka123/GophKeeper/internal/server"
	"github.com/maffka123/GophKeeper/internal/server/config"
	"github.com/maffka123/GophKeeper/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Print("starting...")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	logger, err := basecfg.InitLogger(cfg.Debug, "server")
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger.Info("initializing the service...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize db
	bp := app.GetBasePath()
	db, err := storage.InitDB(ctx, cfg.DBpath, logger, bp)
	if err != nil {
		logger.Fatal("Error initializing db", zap.Error(err))
	}

	srv := server.New(logger, db, cfg.Key)

	authfunc := srv.JWTAuthFunction()
	// run grpc server
	grpcServer := grpc.NewServer(
		// middlewares
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// auth
				grpc_auth.UnaryServerInterceptor(authfunc))),
	)

	reflection.Register(grpcServer)

	pb.RegisterGophKeeperServer(grpcServer, srv)

	listener, err := net.Listen("tcp", cfg.Endpoint)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// handle service stop
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-quit
		logger.Info(fmt.Sprintf("caught sig: %+v", sig))
		grpcServer.GracefulStop()
	}()

	log.Fatal(grpcServer.Serve(listener))
}
