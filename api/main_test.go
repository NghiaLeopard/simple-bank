package api

import (
	"testing"

	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T,store db.Store) *Server {
	testServer,err := NewServer(utils.Config{
		SymmetricKey: utils.RandomString(32),
	},store)

	require.NoError(t,err)

	return testServer
}