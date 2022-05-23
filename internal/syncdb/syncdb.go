package syncdb

import (
	"context"
	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/storage"
	"go.uber.org/zap"
	"time"
)

type SyncDB struct {
	lastSync time.Time
	UserID   int64
	Token    string
	db       storage.StoregeInterface
	c        pb.GophKeeperClient
	ctx      context.Context
	logger   *zap.Logger
}

func NewSyncDB(
	ctx context.Context,
	userID int64,
	token string,
	db storage.StoregeInterface,
	client pb.GophKeeperClient,
	log *zap.Logger) SyncDB {
	return SyncDB{
		UserID: userID,
		Token:  token,
		db:     db,
		c:      client,
		ctx:    ctx,
		logger: log,
	}
}

func (s SyncDB) Sync() error {
	resp, err := s.c.GetAllDataForUser(s.ctx, &pb.GetAllDataForUserRequest{UserID: s.UserID, Time: s.lastSync.String()})
	if err != nil {
		s.logger.Error("data select failed")
		return err
	}

	err = s.db.InserDataForUser(s.ctx, resp.Data, s.UserID)
	s.lastSync = time.Now()

	data, err := s.db.SelectAllDataForUser(s.ctx, s.UserID, s.lastSync.String(), false)
	if err != nil {
		s.logger.Error("data select failed")
		return err
	}
	if len(data) != 0 {
		s.c.InsertSyncData(s.ctx, &pb.InsertSyncDataRequest{Data: data})
	}

	return nil
}

func InitSync(ctx context.Context, tokenChan <-chan string,
	userID <-chan int64,
	db storage.StoregeInterface,
	client pb.GophKeeperClient,
	log *zap.Logger, tsync <-chan time.Time) {

	token := <-tokenChan
	id := <-userID
	s := NewSyncDB(ctx, id, token, db, client, log)
	go s.syncRoutine(tsync)
}

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
