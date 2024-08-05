package token

import (
	"testing"
	"time"

	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := utils.RandomString(32)
	userName := utils.RandomOwner()
	duration := time.Minute * 10

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	pasetoMaker,err := NewPasetoMaker([]byte(symmetricKey))

	require.NoError(t,err)
	require.NotEmpty(t,pasetoMaker)
	
	token,payload, err := pasetoMaker.CreateTokenPaseto(userName,duration)

	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	require.Equal(t,userName,payload.Username)
	require.WithinDuration(t,issuedAt,payload.IssuedAt,time.Second)
	require.WithinDuration(t,expiredAt,payload.Expired,time.Second)

	payload, err = pasetoMaker.VerifyTokenPaseto(token)

	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.Equal(t,userName,payload.Username)
	require.WithinDuration(t,issuedAt,payload.IssuedAt,time.Second)
	require.WithinDuration(t,expiredAt,payload.Expired,time.Second)
}

func TestExpireToken(t *testing.T){
	symmetricKey := utils.RandomString(32)
	userName := utils.RandomOwner()
	
	pasetoMaker,err := NewPasetoMaker([]byte(symmetricKey))

	require.NoError(t,err)
	require.NotEmpty(t,pasetoMaker)
	
	token,payload, err := pasetoMaker.CreateTokenPaseto(userName,-time.Minute)

	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	payload, err = pasetoMaker.VerifyTokenPaseto(token)

	require.Error(t,err)
	require.EqualError(t,err,errExpired.Error())
	require.Empty(t,payload)
}