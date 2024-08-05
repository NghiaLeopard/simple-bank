package db

import (
	"context"
	"testing"
	"time"

	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func CreateUserRandom(t *testing.T) User {
	hashedPassword,err := utils.HashPassword(utils.RandomString(6))

	require.NoError(t,err)
	require.NotEmpty(t,hashedPassword)

	arg := CreateUserParams {
		Username: utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: utils.RandomOwner(),
		Email: utils.RandomEmail(),
	}

	user,err := testQuery.CreateUser(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,user)

	require.Equal(t,user.Email,arg.Email)
	require.Equal(t,hashedPassword,arg.HashedPassword)
	require.Equal(t,user.FullName,arg.FullName)
	require.Equal(t,user.Username,arg.Username)

	require.NotZero(t,user.CreatedAt)
	require.NotZero(t,user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	CreateUserRandom(t)
}

func TestGetUser(t *testing.T) {
	userRandom := CreateUserRandom(t)

	user,err := testQuery.GetUser(context.Background(),userRandom.Username)

	require.NoError(t,err)
	require.NotEmpty(t,user)

	require.Equal(t,userRandom.Email,user.Email)
	require.Equal(t,userRandom.FullName,user.FullName)
	require.Equal(t,userRandom.HashedPassword,user.HashedPassword)
	require.Equal(t,userRandom.Username,user.Username)

	require.WithinDuration(t,userRandom.CreatedAt,user.CreatedAt,time.Second)
}

