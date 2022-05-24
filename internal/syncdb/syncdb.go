package syncdb

import (
	"context"
	"fmt"
	"time"

	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SyncDB struct {
	lastSync string
	UserID   int64
	Token    string
	db       storage.StoregeInterface
	c        pb.GophKeeperClient
	ctx      context.Context
	logger   *zap.Logger
}

// NewSyncDB returns new sync db object
func NewSyncDB(
	ctx context.Context,
	userID int64,
	token string,
	db storage.StoregeInterface,
	client pb.GophKeeperClient,
	log *zap.Logger) SyncDB {
	return SyncDB{
		lastSync: time.Now().Format("2006-01-02 15:04:05"), // TODO: implement proper time extraction (last synced time from local db)
		UserID:   userID,
		Token:    token,
		db:       db,
		c:        client,
		ctx:      ctx,
		logger:   log,
	}
}

// Sync synchronizes dbs
// TODO: sync users tables
func (s SyncDB) Sync() error {
	s.ctx = metadata.AppendToOutgoingContext(s.ctx, "token", s.Token)
	resp, err := s.c.GetAllDataForUser(s.ctx, &pb.GetAllDataForUserRequest{UserID: s.UserID, Time: s.lastSync})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Code(codes.Unavailable):
				s.logger.Warn("server is not available.....")
				return nil
			}
		} else {
			s.logger.Error(fmt.Sprintf("data select from server failed: %s", err.Error()))
			return err
		}
	}

	err = s.db.InserDataForUser(s.ctx, resp.Data, s.UserID)
	s.lastSync = time.Now().Format("2006-01-02 15:04:05")

	data, err := s.db.SelectAllDataForUser(s.ctx, s.UserID, s.lastSync, false)
	if err != nil {
		s.logger.Error(fmt.Sprintf("data select from client failed: %s", err.Error()))
		return err
	}
	if len(data) != 0 {
		_, err = s.c.InsertSyncData(s.ctx, &pb.InsertSyncDataRequest{Data: data})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Code(codes.Unavailable):
					s.logger.Warn("server is not available.....")
					return nil
				}
			} else {
				s.logger.Error(fmt.Sprintf("data insert to server failed: %s", err.Error()))
				return err
			}
		}
	}

	return nil
}

// InitSync starts synchronizing dbs as soon as it has user id and token
// TODO: problem with multiple users
func InitSync(ctx context.Context, tokenChan chan string, userID chan int64,
	db storage.StoregeInterface, client pb.GophKeeperClient, log *zap.Logger, tsync <-chan time.Time) {

	token := <-tokenChan
	id := <-userID
	s := NewSyncDB(ctx, id, token, db, client, log)
	go s.syncRoutine(tsync)
}

// syncRoutine runs periodic sync of dbs as go routine
func (s SyncDB) syncRoutine(t <-chan time.Time) {
	for {
		select {
		case <-t:
			s.logger.Info("synchronizing dbs")
			s.Sync()
		case <-s.ctx.Done():
			s.logger.Info("context canceled")
		}
	}
}
