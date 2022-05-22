package storage

import (
	"testing"

	rpc "github.com/maffka123/GophKeeper/api/proto"
	"github.com/stretchr/testify/assert"
)

func TestPGDB_buildQuery(t *testing.T) {
	db := &PGDB{}

	tests := []struct {
		name string
		data *rpc.Data
		want string
	}{
		{name: "with id",
			data: &rpc.Data{ID: 1, Data: &rpc.KeepData{AuthData: &rpc.AuthData{Login: "name"}}},
			want: "SELECT * FROM secrets WHERE user_id=0 AND id=1",
		},

		{name: "with meta and login",
			data: &rpc.Data{Data: &rpc.KeepData{AuthData: &rpc.AuthData{Login: "name"}}, Metadata: "some data"},
			want: `SELECT * FROM secrets WHERE user_id=0 AND metadata LIKE '%some data%' AND data #> '{AuthData}' @? '$.Login ? (@ == "name")'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, _ := db.buildQuery(tt.data, "*")

			assert.Equal(t, s, tt.want)
		})
	}
}
