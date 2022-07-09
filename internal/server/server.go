package server

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/app"
	"github.com/maffka123/GophKeeper/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// secretService struct for grpc service
type secretService struct {
	pb.UnimplementedGophKeeperServer
	logger *zap.Logger
	db     storage.StoregeInterface
	token  *jwtauth.JWTAuth
}

// New creates new instance of grpc service
func New(logger *zap.Logger, db storage.StoregeInterface, secret string) *secretService {
	return &secretService{
		logger: logger,
		db:     db,
		token:  jwtauth.New("HS256", []byte(secret), nil),
	}
}

// Register registers new user
func (s *secretService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResp, error) {
	exists, err := s.db.CreateNewUser(ctx, request.User)
	if exists == -1 {
		return &pb.RegisterResp{Message: "User already exists", Exists: -1}, nil
	}
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	_, tokenString, _ := s.token.Encode(map[string]interface{}{"user_id": request.User.ID})
	return &pb.RegisterResp{Message: "Successfully created user", Exists: 0, Token: tokenString}, nil
}

// Login checks if given password is correct and issues authorisation token
func (s *secretService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResp, error) {
	pass, id, err := s.db.SelectPass(ctx, request.User)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	if pass == nil || !comparePass(*pass, request.User.Password) {
		return nil,
			status.Errorf(
				codes.PermissionDenied, "login or password are wrong",
			)
	}

	_, tokenString, _ := s.token.Encode(map[string]interface{}{"user_id": request.User.ID})
	return &pb.LoginResp{Message: "Logged in successfully", Token: tokenString, UserId: *id}, nil
}

// Insert inserts data to postgres from authorized user
func (s *secretService) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResp, error) {
	currUser, err := app.UserIDFromContext(ctx)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	s.logger.Debug("found user: ", zap.String("login", fmt.Sprint(currUser)))

	request.Data.UserID = currUser
	id, err := s.db.InsertData(ctx, request.Data, true)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	return &pb.InsertResp{Id: fmt.Sprint(id)}, nil
}

// GetData gets data for athorized user
func (s *secretService) GetData(ctx context.Context, request *pb.GetDataRequest) (*pb.GetDataResp, error) {
	currUser, err := app.UserIDFromContext(ctx)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	s.logger.Debug("found user: ", zap.String("login", fmt.Sprint(currUser)))

	request.Data.UserID = currUser
	data, err := s.db.SearchData(ctx, request.Data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	return &pb.GetDataResp{Data: data}, nil
}

// Delete deletes data for athorized user
func (s *secretService) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResp, error) {
	currUser, err := app.UserIDFromContext(ctx)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	s.logger.Debug("found user: ", zap.String("login", fmt.Sprint(currUser)))

	request.Data.UserID = currUser
	data, err := s.db.DeleteData(ctx, request.Data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	return &pb.DeleteResp{Data: data}, nil
}

func (s *secretService) GetAllDataForUser(ctx context.Context, request *pb.GetAllDataForUserRequest) (*pb.GetAllDataForUserResp, error) {
	currUser, err := app.UserIDFromContext(ctx)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	s.logger.Debug("found user: ", zap.String("login", fmt.Sprint(currUser)))

	data, err := s.db.SelectAllDataForUser(ctx, currUser, request.Time, true)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	return &pb.GetAllDataForUserResp{Data: data}, nil
}

func (s *secretService) InsertSyncData(ctx context.Context, request *pb.InsertSyncDataRequest) (*pb.InsertSyncDataResp, error) {

	currUser, err := app.UserIDFromContext(ctx)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	s.logger.Debug("found user: ", zap.String("login", fmt.Sprint(currUser)))

	err = s.db.InserDataForUser(ctx, request.Data, currUser)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, err.Error(),
		)
	}
	return &pb.InsertSyncDataResp{Message: "data inserted"}, nil
}

// ComparePass compares hashed passwords.
func comparePass(expected string, actual string) bool {
	np := sha256.Sum256([]byte(actual))
	npb := hex.EncodeToString(np[:])
	return subtle.ConstantTimeCompare([]byte(expected), []byte(npb)) == 1
}

// JWTAuthFunction generated func of type jwtAuth to use as a middleware for grpc server
func (s *secretService) JWTAuthFunction() func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {

		// ignore some endpoints for jwt check
		method, _ := grpc.Method(ctx)
		for _, imethod := range []string{"/proto.GophKeeper/Register", "/proto.GophKeeper/Login"} {
			if method == imethod {
				return ctx, nil
			}
		}

		tokenString := metautils.ExtractIncoming(ctx).Get("token")
		token, err := jwtauth.VerifyToken(s.token, tokenString)

		if err != nil {
			return nil, status.Errorf(
				codes.Unauthenticated, err.Error(),
			)
		}
		ctx = jwtauth.NewContext(ctx, token, err)
		return ctx, nil
	}
}
