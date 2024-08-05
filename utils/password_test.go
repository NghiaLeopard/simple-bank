package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// when write file test , must write case wrong case
func TestPassword(t *testing.T) {
	password := RandomString(8)

	newPassword,err := HashPassword(password)

	require.NoError(t,err)

	require.NotEmpty(t,newPassword)

	err = CheckPassword(password,newPassword)
	t.Log(password,newPassword)


	require.NoError(t,err)

	wrongPassword := RandomString(8)

	err = CheckPassword(wrongPassword,newPassword)

	require.EqualError(t,err,bcrypt.ErrMismatchedHashAndPassword.Error())
}