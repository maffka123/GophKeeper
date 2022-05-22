package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	rpc "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/server/config"
	"go.uber.org/zap"
	//"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// PGinterface interface of postgres data sources
type PGinterface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Close()
}

// StoregeInterface interface for storage
type StoregeInterface interface {
	CreateNewUser(context.Context, *rpc.User) (int, error)
	SelectPass(context.Context, *rpc.User) (*string, error)
	InsertData(context.Context, *rpc.Data) (*int64, error)
	SearchData(context.Context, *rpc.Data) ([]*rpc.Data, error)
	DeleteData(context.Context, *rpc.Data) ([]*rpc.Data, error)
}

// PGDB type for postgres
type PGDB struct {
	path string
	Conn PGinterface
	log  *zap.Logger
}

// InitDB initialized pg connection and creates tables
func InitDB(ctx context.Context, cfg *config.Config, logger *zap.Logger, bp string) (*PGDB, error) {
	db := PGDB{
		path: cfg.DBpath,
		log:  logger,
	}
	conn, err := pgxpool.Connect(ctx, cfg.DBpath)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	db.Conn = conn

	db.log.Info("initializing db tables...")
	file, err := ioutil.ReadFile(filepath.Join(bp, "/docker_db/db.sql"))
	if err != nil {
		return nil, fmt.Errorf("unable to read sql file: %v", err)
	}

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %v", err)
	}
	defer tx.Rollback(ctx)

	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(ctx, q); err != nil {
			return nil, fmt.Errorf("failed executing sql: %v", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit failed: %v", err)
	}
	db.log.Info("db initialized succesfully")

	return &db, nil
}

// CreateNewUser insertes new user, handles not unique users
func (db *PGDB) CreateNewUser(ctx context.Context, user *rpc.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	np := sha256.Sum256([]byte(user.Password))
	pass := hex.EncodeToString((np[:]))

	err := db.Conn.QueryRow(ctx, `INSERT INTO users (login, password) VALUES($1,$2) RETURNING id;`, user.Login, pass).Scan(&user.ID)

	if err != nil && strings.Contains(err.Error(), "violates") {
		return -1, fmt.Errorf("user already exists: %v", err)
	} else if err != nil {
		return 0, fmt.Errorf("insert new user failed: %v", err)
	}

	return 1, nil
}

// SelectPass gets hashed password for a particular user
func (db *PGDB) SelectPass(ctx context.Context, user *rpc.User) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var val string
	row := db.Conn.QueryRow(ctx, "SELECT password, id FROM users WHERE login=$1", user.Login)
	err := row.Scan(&val, &user.ID)

	if err != nil {
		return nil, fmt.Errorf("select from users failed: %v", err)
	}
	// TODO user does not exist
	return &val, nil
}

// InsertData appends new data to existing secrets
func (db *PGDB) InsertData(ctx context.Context, data *rpc.Data) (*int64, error) {
	var id int64
	err := db.Conn.QueryRow(ctx, `
	INSERT INTO secrets (user_id, data, metadata) VALUES ($1,$2,$3) RETURNING id`,
		data.UserID, data.Data, data.Metadata).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// SearchData gets all secrets that match search criteria
func (db *PGDB) SearchData(ctx context.Context, in *rpc.Data) ([]*rpc.Data, error) {
	var out []*rpc.Data
	query, err := db.buildQuery(in, "id, data, metadata")
	if err != nil {
		return nil, err
	}

	row, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("select from secrets failed: %v", err)
	}
	defer row.Close()

	for row.Next() {
		var o rpc.Data
		var g rpc.KeepData
		err := row.Scan(&o.ID, &g, &o.Metadata)
		o.Data = &g
		if err != nil {
			db.log.Error("select from secrets failed:", zap.Error(err))
		}

		if err != nil {
			db.log.Error("unmarshar data from secrets failed:", zap.Error(err))
		}
		out = append(out, &o)
	}

	return out, nil
}

// DeleteData gets all secrets that match search criteria
func (db *PGDB) DeleteData(ctx context.Context, in *rpc.Data) ([]*rpc.Data, error) {
	var ids []string
	data, err := db.SearchData(ctx, in)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		ids = append(ids, fmt.Sprint(d.ID))
	}

	_, err = db.Conn.Exec(ctx, `DELETE FROM secrets WHERE id in ($1);`, strings.Join(ids, ","))

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *PGDB) buildQuery(data *rpc.Data, cols string) (string, error) {
	query := fmt.Sprintf("SELECT %s FROM secrets WHERE user_id=%d AND ", cols, data.UserID)
	if data.ID != 0 {
		return query + fmt.Sprintf("id=%d", data.ID), nil
	}

	var filter []string
	if data.Metadata != "" {
		filter = append(filter, fmt.Sprintf("metadata LIKE '%%%s%%'", data.Metadata))
	}

	if data.Data.GetAuthData() != nil && data.Data.AuthData.Login != "" {
		filter = append(filter, fmt.Sprintf(`data #> '{AuthData}' @? '$.login ? (@ == "%s")'`, data.Data.AuthData.Login))
	}

	if data.Data.GetText() != "" {
		filter = append(filter, fmt.Sprintf(`data @? '$.Text ? (@ == "%s")'`, data.Data.Text))
	}

	if data.Data.GetBinary() != nil {
		b := string(data.Data.Binary[:])
		filter = append(filter, fmt.Sprintf(`data @? '$.Binary ? (@ == "%s")'`, b))
	}

	if data.Data.GetBankCard() != nil {
		if data.Data.BankCard.Address != "" {
			filter = append(filter, fmt.Sprintf(`data #> '{BankCard}' @? '$.Address ? (@ == "%s")'`, data.Data.BankCard.Address))
		}

		if data.Data.BankCard.BankName != "" {
			filter = append(filter, fmt.Sprintf(`data #> '{BankCard}' @? '$.BankName ? (@ == "%s")'`, data.Data.BankCard.BankName))
		}

		if data.Data.BankCard.Expiry != "" {
			filter = append(filter, fmt.Sprintf(`data #> '{BankCard}' @? '$.Expiry ? (@ == "%s")'`, data.Data.BankCard.Expiry))
		}

		if data.Data.BankCard.HolderName != "" {
			filter = append(filter, fmt.Sprintf(`data #> '{BankCard}' @? '$.HolderName ? (@ == "%s")'`, data.Data.BankCard.HolderName))
		}

		if data.Data.BankCard.CardNumber != 0 {
			filter = append(filter, fmt.Sprintf(`data #> '{BankCard}' @? '$.CardNumber ? (@ == "%d")'`, data.Data.BankCard.CardNumber))
		}
	}

	if len(filter) > 0 {
		query += strings.Join(filter[:], " AND ")
		return query, nil
	}
	return "", fmt.Errorf("could not construct any filter, please add more data")
}
