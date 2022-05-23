package main

import (
	"log"

	"context"
	"fmt"
	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/app"
	"github.com/maffka123/GophKeeper/internal/client/config"
	basecfg "github.com/maffka123/GophKeeper/internal/config"
	"github.com/maffka123/GophKeeper/internal/handlers"
	"github.com/maffka123/GophKeeper/internal/storage"
	"github.com/maffka123/GophKeeper/internal/syncdb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Version     string = "N/A"
	BuildDate   string = "N/A"
	BuildCommit string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", Version, BuildDate, BuildCommit)
	fmt.Print("starting...")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	logger, err := basecfg.InitLogger(cfg.Debug, cfg.AppName)
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

	conn, err := grpc.Dial(cfg.ServerEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewGophKeeperClient(conn)

	defer conn.Close()

	// initialize sync
	tokenChan := make(chan string)
	idChan := make(chan int64)
	syncTicker := time.NewTicker(cfg.SyncInterval)
	go syncdb.InitSync(ctx, tokenChan, idChan, db, client, logger, syncTicker.C)

	// prepare handles
	r := handlers.KeeperRouter(ctx, logger, client, db, tokenChan, idChan)

	// handle service stop
	srv := &http.Server{Addr: cfg.Endpoint, Handler: r}
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-quit
		logger.Info(fmt.Sprintf("caught sig: %+v", sig))
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Error("HTTP server Shutdown:", zap.Error(err))
		}
	}()

	logger.Info("Start serving on", zap.String("endpoint name", cfg.Endpoint))
	log.Fatal(srv.ListenAndServe())

}
