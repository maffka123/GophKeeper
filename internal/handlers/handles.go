package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	pb "github.com/maffka123/GophKeeper/api/proto"
	"github.com/maffka123/GophKeeper/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// Handler struct for api handler
type Handler struct {
	logger *zap.Logger
	ctx    context.Context
	c      pb.GophKeeperClient
	db     storage.StoregeInterface
}

// NewHandler returns new initilized handler
func NewHandler(ctx context.Context, logger *zap.Logger, c pb.GophKeeperClient, db storage.StoregeInterface) Handler {
	return Handler{
		logger: logger,
		ctx:    ctx,
		c:      c,
		db:     db,
	}
}

// HandlerPostRegister creates new user if user with such login not yet exist
func (h *Handler) HandlerPostRegister(tokenChan chan string, idChan chan int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Register json cannot be read: %s", err), http.StatusBadRequest)
			return
		}
		var u pb.User
		err = protojson.Unmarshal(body, &u)
		h.logger.Debug("recieved new user: ", zap.String("login", u.Login))
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Register json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}

		resp, err := h.c.Register(h.ctx, &pb.RegisterRequest{User: &u})
		if resp.Exists == -1 {
			http.Error(w, fmt.Sprintf("409 - Login is already taken: %s", err), http.StatusConflict)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
			return
		} else if resp.Exists == 0 {

			h.logger.Debug("logged in: ", zap.String("login", u.Login))
			http.SetCookie(w, &http.Cookie{
				Name:  "jwt",
				Value: resp.Token,
			})
			tokenChan <- resp.Token
			idChan <- resp.Exists

			w.Header().Set("application-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok}`))
		}
	}
}

// HandlerPostLogin logins user if login and password are valid
func (h *Handler) HandlerPostLogin(tokenChan chan string, idChan chan int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Login json cannot be read: %s", err), http.StatusBadRequest)
			return
		}
		var u pb.User
		err = protojson.Unmarshal(body, &u)
		h.logger.Debug("user is trying to login: ", zap.String("login", u.Login))

		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Login json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}

		resp, err := h.c.Login(h.ctx, &pb.LoginRequest{User: &u})

		if err != nil {
			http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
			return
		}

		h.logger.Debug("logged in: ", zap.String("login", u.Login))
		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: resp.Token,
		})

		tokenChan <- resp.Token
		idChan <- resp.UserId

		w.Header().Set("application-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok}`))

	}
}

// HandlerPostData adds new secret for current user
func (h *Handler) HandlerPostData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Data json cannot be read: %s", err), http.StatusBadRequest)
			return
		}
		var d pb.Data
		err = protojson.Unmarshal(body, &d)

		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Data json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}
		tokenstr := jwtauth.TokenFromCookie(r)
		h.ctx = metadata.AppendToOutgoingContext(h.ctx, "token", tokenstr)

		resp, err := h.c.Insert(h.ctx, &pb.InsertRequest{Data: &d})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Code(codes.Unavailable):
					h.logger.Warn("server is not available saving data locally")
					h.insertIfUnavailable(d)
				}
			} else {
				h.logger.Debug(err.Error())
				http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
				return
			}
		}

		h.logger.Debug("order accepted: ", zap.String("login", fmt.Sprint(resp.Id)))
		w.Header().Set("application-type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"status":"ok, "insert id": %d}`, resp.Id)))
	}
}

func (h *Handler) insertIfUnavailable(data pb.Data) (*string, error) {
	id, err := h.db.InsertData(h.ctx, &data, false)
	if err != nil {
		return nil, err
	}
	idstr := fmt.Sprint(id)
	return &idstr, nil
}

// HandlerGetData gets data for current user for given search criteria
func (h *Handler) HandlerGetData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Data json cannot be read: %s", err), http.StatusBadRequest)
			return
		}
		var d pb.Data
		err = protojson.Unmarshal(body, &d)

		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Data json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}

		tokenstr := jwtauth.TokenFromCookie(r)
		h.ctx = metadata.AppendToOutgoingContext(h.ctx, "token", tokenstr)

		resp, err := h.c.GetData(h.ctx, &pb.GetDataRequest{Data: &d})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Code(codes.Unavailable):
					h.logger.Warn("server is not available getting data from local base")
					h.getDataIfUnavailable(d)
				}
			} else {
				h.logger.Debug(err.Error())
				http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
				return
			}
		}

		data, err := protojson.Marshal(resp)
		if err != nil {
			h.logger.Debug(err.Error())
			http.Error(w, fmt.Sprintf("500 - could not convert data to json: %s", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (h *Handler) getDataIfUnavailable(data pb.Data) ([]*pb.Data, error) {
	d, err := h.db.SearchData(h.ctx, &data)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// HandlerGetDelete deletes for current user for given search criteria
func (h *Handler) HandlerGetDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Data json cannot be read: %s", err), http.StatusBadRequest)
			return
		}
		var d pb.Data
		err = protojson.Unmarshal(body, &d)

		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Data json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}

		tokenstr := jwtauth.TokenFromCookie(r)
		h.ctx = metadata.AppendToOutgoingContext(h.ctx, "token", tokenstr)

		resp, err := h.c.Delete(h.ctx, &pb.DeleteRequest{Data: &d})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Code(codes.Unavailable):
					h.logger.Warn("server is not available getting data from local base")
					h.deleteDataIfUnavailable(d)
				}
			} else {
				h.logger.Debug(err.Error())
				http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
				return
			}
		}

		data, err := protojson.Marshal(resp)
		if err != nil {
			h.logger.Debug(err.Error())
			http.Error(w, fmt.Sprintf("500 - could not convert data to json: %s", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (h *Handler) deleteDataIfUnavailable(data pb.Data) ([]*pb.Data, error) {
	d, err := h.db.DeleteData(h.ctx, &data)
	if err != nil {
		return nil, err
	}
	return d, nil
}
