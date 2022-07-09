package app

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/go-chi/jwtauth/v5"
)

// GetBasePath prepares base path of the project.
func GetBasePath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(path.Dir(b)))
}

// UserIDFromContext gets user if exists from context (it gets there from jwt token parsing).
func UserIDFromContext(ctx context.Context) (int64, error) {
	_, uID, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	if id, ok := uID["user_id"].(float64); ok {
		return int64(id), nil
	}

	return 0, fmt.Errorf("user_id could not be parsed a number: %v", uID["user_id"])
}
