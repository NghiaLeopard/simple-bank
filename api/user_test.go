package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockstore "github.com/NghiaLeopard/simple-bank/db/mock"
	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type req struct{
	Username       string `json:"username"`
	Password 	   string `json:"password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
	password := utils.RandomString(6)

	hashedPassword,err := utils.HashPassword(password)

	require.NoError(t,err)

	randomUser :=  db.User{
		Username: utils.RandomOwner(),
		FullName: utils.RandomOwner(),
		HashedPassword: hashedPassword,
		Email: utils.RandomEmail(),
	}

	arg := db.CreateUserParams{
		Username: randomUser.Username,
		FullName: randomUser.FullName,
		Email: randomUser.Email,
	}

	testcase := []struct{
		name string
		body req
		buildStubs func(store *mockstore.MockStore)
		checkResponse func(t *testing.T,recorder * httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: req{
					Username: randomUser.Username,
					FullName: randomUser.FullName,
					Password: password,
					Email: randomUser.Email,
			},
			buildStubs: func(store *mockstore.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),EqCreateUserParams(arg,password)).Times(1).Return(randomUser,nil)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t,http.StatusOK,recorder.Code)
				requiredMatchBody(t,recorder.Body,randomUser)
			},
		},
		{
			name: "internal server error",
			body: req{
					Username: randomUser.Username,
					FullName: randomUser.FullName,
					Password: password,
					Email: randomUser.Email,
			},
			buildStubs: func(store *mockstore.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(1).Return(db.User{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder) {
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
		{
			name: "Too short password",
			body: req{
				Username: randomUser.Username,
				FullName: randomUser.FullName,
				Password: "sadsd",
				Email: randomUser.Email,
		},
			buildStubs: func(store *mockstore.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder) {
				require.Equal(t,recorder.Code,http.StatusBadRequest)
			},
		},
		{
			name: "Wrong email",
			body: req{
				Username: "nghia@1 ",
				FullName: randomUser.FullName,
				Password: password,
				Email: randomUser.Email,
		},
			buildStubs: func(store *mockstore.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder) {
				require.Equal(t,recorder.Code,http.StatusBadRequest)
			},
		},
		{
			name: "Invalid username",
			body: req{
				Username: randomUser.Username,
				FullName: randomUser.FullName,
				Password: "sadsd",
				Email: "nghiabeo1605gmail.com",
		},
			buildStubs: func(store *mockstore.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder) {
				require.Equal(t,recorder.Code,http.StatusBadRequest)
			},
		},
	}

	for i:= range testcase {
		tc := testcase[i]

		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store :=  mockstore.NewMockStore(ctrl)

			tc.buildStubs(store)
		
			server := newTestServer(t,store)
		
			recorder := httptest.NewRecorder()

		
			dataParse,err := utils.JsonReaderFactory(tc.body)
		
			require.NoError(t,err)

			t.Log(dataParse)

			request,err := http.NewRequest(http.MethodPost,"/user",dataParse)
		
			require.NoError(t,err)
		
			server.router.ServeHTTP(recorder,request)

			tc.checkResponse(t,recorder)
		
		})
	}
	
}


func requiredMatchBody(t *testing.T,body *bytes.Buffer,account db.User){
	data,err := io.ReadAll(body)

	require.NoError(t,err)

	var gotAccount db.User

	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t,err)

	require.Equal(t,gotAccount,account)
}