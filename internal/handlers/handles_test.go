package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/protobuf/encoding/protojson"

	"log"
	"net"

	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/server"
	"github.com/maffka123/GophKeeper/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestHandler_HandlerPostRegister(t *testing.T) {
	ctx := context.Background()
	//init stuff
	logger, _ := zap.NewDevelopment()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGophKeeperClient(conn)
	r := KeeperRouter(context.Background(), logger, client)

	type want struct {
		statusCode int
	}
	type request struct {
		route string
		body  pb.User
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{name: "register_success",
			request: request{route: "/api/user/register", body: pb.User{Login: "test", Password: "pass"}},
			want:    want{statusCode: 200},
		},
		{name: "user_exists",
			request: request{route: "/api/user/register", body: pb.User{Login: "error", Password: "pass"}},
			want:    want{statusCode: 409},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, _ := json.Marshal(tt.request.body)

			request := httptest.NewRequest(http.MethodPost, tt.request.route, bytes.NewBuffer(body))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_HandlerPostLogin(t *testing.T) {
	ctx := context.Background()
	//init stuff
	logger, _ := zap.NewDevelopment()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGophKeeperClient(conn)
	r := KeeperRouter(context.Background(), logger, client)
	type want struct {
		statusCode int
		cookie     http.Cookie
	}
	type request struct {
		route string
		body  pb.User
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{name: "login_success",
			request: request{route: "/api/user/login", body: pb.User{Login: "test", Password: "pass"}},
			want:    want{statusCode: 200, cookie: http.Cookie{Name: "jwt", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMX0.l7nSrATz05XodnENMEpOZfvDLK3FzkbuzNWdliJ7Alo"}},
		},
		{name: "user_not_exist",
			request: request{route: "/api/user/login", body: pb.User{Login: "error", Password: "pass"}},
			want:    want{statusCode: 500},
		},
		{name: "pass_wrong",
			request: request{route: "/api/user/login", body: pb.User{Login: "test", Password: "pass1"}},
			want:    want{statusCode: 500},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, _ := json.Marshal(tt.request.body)

			request := httptest.NewRequest(http.MethodPost, tt.request.route, bytes.NewBuffer(body))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.statusCode == 200 {
				assert.Equal(t, tt.want.cookie.Value, result.Cookies()[0].Value)
				assert.Equal(t, tt.want.cookie.Name, result.Cookies()[0].Name)
			}

		})
	}
}

func TestHandler_HandlerPostData(t *testing.T) {
	ctx := context.Background()
	//init stuff
	logger, _ := zap.NewDevelopment()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGophKeeperClient(conn)
	r := KeeperRouter(context.Background(), logger, client)
	type want struct {
		statusCode int
	}
	type request struct {
		route string
		body  pb.Data
	}
	tests := []struct {
		name    string
		request request
		want    want
		db      fakeDB
	}{
		{name: "data_added",
			request: request{route: "/api/user/insert", body: pb.Data{Data: &pb.KeepData{AuthData: &pb.AuthData{Login: "name", Password: "pass"}}}},
			want:    want{statusCode: 202},
			db:      fakeDB{selectUserForOrder: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, _ := protojson.Marshal(&tt.request.body)

			request := httptest.NewRequest(http.MethodPost, tt.request.route, bytes.NewBuffer(body))
			request.AddCookie(&http.Cookie{Name: "jwt", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMX0.l7nSrATz05XodnENMEpOZfvDLK3FzkbuzNWdliJ7Alo"})
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

		})
	}
}

func TestHandler_HandlerGetOrders(t *testing.T) {
	ctx := context.Background()
	//init stuff
	logger, _ := zap.NewDevelopment()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGophKeeperClient(conn)
	r := KeeperRouter(context.Background(), logger, client)
	type want struct {
		statusCode int
	}
	type request struct {
		route string
		body  pb.Data
	}
	tests := []struct {
		name    string
		request request
		want    want
		db      fakeDB
	}{
		{name: "empty_orders",
			request: request{route: "/api/user/search", body: pb.Data{Data: &pb.KeepData{AuthData: &pb.AuthData{Login: "name"}}}},
			want:    want{statusCode: 202},
			db:      fakeDB{SearchDataRes: []*pb.Data{}},
		},
		{name: "get_orders",
			request: request{route: "/api/user/search", body: pb.Data{ID: 1}},
			want:    want{statusCode: 202},
			db: fakeDB{SearchDataRes: []*pb.Data{{ID: 1, Data: &pb.KeepData{Text: "smth1"}},
				{ID: 2, Data: &pb.KeepData{Text: "smth2"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := protojson.Marshal(&tt.request.body)
			request := httptest.NewRequest(http.MethodGet, tt.request.route, bytes.NewBuffer(body))
			request.AddCookie(&http.Cookie{Name: "jwt", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMX0.l7nSrATz05XodnENMEpOZfvDLK3FzkbuzNWdliJ7Alo"})
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.statusCode == 200 {
				var d *pb.GetDataResp
				body, _ := ioutil.ReadAll(result.Body)
				protojson.Unmarshal(body, d)
				assert.Equal(t, d.Data, tt.db.SearchDataRes)
			}

		})
	}
}

func TestHandler_HandlerGetDelete(t *testing.T) {
	ctx := context.Background()
	//init stuff
	logger, _ := zap.NewDevelopment()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGophKeeperClient(conn)
	r := KeeperRouter(context.Background(), logger, client)
	type want struct {
		statusCode int
	}
	type request struct {
		route string
		body  pb.Data
	}
	tests := []struct {
		name    string
		request request
		want    want
		db      fakeDB
	}{
		{name: "empty_orders",
			request: request{route: "/api/user/search", body: pb.Data{Data: &pb.KeepData{AuthData: &pb.AuthData{Login: "name"}}}},
			want:    want{statusCode: 202},
			db:      fakeDB{SearchDataRes: []*pb.Data{}},
		},
		{name: "get_orders",
			request: request{route: "/api/user/search", body: pb.Data{ID: 1}},
			want:    want{statusCode: 202},
			db: fakeDB{SearchDataRes: []*pb.Data{{ID: 1, Data: &pb.KeepData{Text: "smth1"}},
				{ID: 2, Data: &pb.KeepData{Text: "smth2"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := protojson.Marshal(&tt.request.body)
			request := httptest.NewRequest(http.MethodGet, tt.request.route, bytes.NewBuffer(body))
			request.AddCookie(&http.Cookie{Name: "jwt", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMX0.l7nSrATz05XodnENMEpOZfvDLK3FzkbuzNWdliJ7Alo"})
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.statusCode == 200 {
				var d *pb.GetDataResp
				body, _ := ioutil.ReadAll(result.Body)
				protojson.Unmarshal(body, d)
				assert.Equal(t, d.Data, tt.db.SearchDataRes)
			}

		})
	}
}

//http://www.inanzzz.com/index.php/post/w9qr/unit-testing-golang-grpc-client-and-server-application-with-bufconn-package
func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	logger, _ := zap.NewDevelopmentConfig().Build()
	db := newFakeDB()
	mysrv := server.New(logger, db, "secret")

	authfunc := mysrv.JWTAuthFunction()
	srv := grpc.NewServer(
		// middlewares
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// auth
				grpc_auth.UnaryServerInterceptor(authfunc))),
	)

	pb.RegisterGophKeeperServer(srv, mysrv)

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type fakeDB struct {
	selectUserForOrder int64
	Conn               storage.PGinterface
	SearchDataRes      []*pb.Data
	DeleteDataRes      []*pb.Data
}

func newFakeDB() *fakeDB {
	return &fakeDB{}
}

func (db *fakeDB) CreateNewUser(ctx context.Context, user *pb.User) (int, error) {
	if user.Login == "error" {
		return -1, fmt.Errorf("user already exists")
	}
	return 1, nil
}

func (db *fakeDB) SelectPass(ctx context.Context, user *pb.User) (*string, error) {
	if user.Login == "error" {
		return nil, fmt.Errorf("user not found")
	}
	np := sha256.Sum256([]byte("pass"))
	npb := hex.EncodeToString(np[:])
	user.ID = 11
	return &npb, nil
}

func (db *fakeDB) SelectUserForOrder(ctx context.Context, d *pb.Data) (int64, error) {
	return db.selectUserForOrder, nil
}
func (db *fakeDB) InsertData(ctx context.Context, d *pb.Data) (*int64, error) {
	i := int64(1)
	return &i, nil
}

func (db *fakeDB) SearchData(ctx context.Context, d *pb.Data) ([]*pb.Data, error) {
	return db.SearchDataRes, nil
}

func (db *fakeDB) DeleteData(ctx context.Context, d *pb.Data) ([]*pb.Data, error) {
	return db.DeleteDataRes, nil
}
