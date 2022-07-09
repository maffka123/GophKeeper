package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/storage"
	"go.uber.org/zap"
)

// KeeperRouter arranges the whole API endpoints and their correponding handlers
func KeeperRouter(ctx context.Context, logger *zap.Logger, c pb.GophKeeperClient,
	db storage.StoregeInterface, tokenChan chan string, idChan chan int64) chi.Router {

	r := chi.NewRouter()
	mh := NewHandler(ctx, logger, c, db)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", Conveyor(mh.HandlerPostRegister(tokenChan, idChan), unpackGZIP, checkForJSON))
		r.Post("/login", Conveyor(mh.HandlerPostLogin(tokenChan, idChan), unpackGZIP, checkForJSON))
		r.Post("/insert", Conveyor(mh.HandlerPostData(), unpackGZIP, checkForJSON))
		r.Get("/search", Conveyor(mh.HandlerGetData(), unpackGZIP, packGZIP))
		r.Get("/delete", Conveyor(mh.HandlerGetDelete(), unpackGZIP, packGZIP))

	})

	return r
}
