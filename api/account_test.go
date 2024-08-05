package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockstore "github.com/NghiaLeopard/simple-bank/db/mock"
	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	accountRandom := RandomAccount()

	testcase := []struct{
		name string
		accountRandom db.Account
		buildStubs func(store *mockstore.MockStore)
		checkResponse func(t *testing.T,recorder * httptest.ResponseRecorder)

	}{
	{	name: "OK",
		buildStubs: func(store *mockstore.MockStore){
			store.EXPECT().CreateAccount(gomock.Any(),gomock.Eq(accountRandom)).Times(1).Return(accountRandom,nil)
		},
		checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder){
			require.Equal(t,recorder.Code,http.StatusOK)
			requiredBodyMatchAccount(t,recorder.Body,accountRandom)
		},
	},
		{name: "Internal server",
		buildStubs: func(store *mockstore.MockStore){
			account := db.CreateAccountParams{
				Owner: accountRandom.Owner,
				Balance: 0,
				Currency: accountRandom.Currency,
			}

			store.EXPECT().CreateAccount(gomock.Any(),gomock.Eq(account)).Times(1).Return(db.Account{},sql.ErrConnDone)
		},
		checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder){
			require.Equal(t,recorder.Code,http.StatusInternalServerError)
		},
	},
	}

	for i := range testcase {
		tc := testcase[i]

		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockstore.NewMockStore(ctrl)
	
			tc.buildStubs(store)
	
			server := newTestServer(t,store)
		
			recorder := httptest.NewRecorder()
		
			dataParse,err := utils.JsonReaderFactory(struct{
				Owner string `json:"owner"`
				Currency string `json:"currency"`
			}{
				Owner: accountRandom.Owner ,
				Currency: accountRandom.Currency,
			})
		
			if err != nil {
				require.Equal(t,recorder.Code,http.StatusBadRequest)
			}
		
			request,err := http.NewRequest(http.MethodPost,"/account",dataParse)
	
			require.NoError(t,err)
		
			server.router.ServeHTTP(recorder,request)
	
			tc.checkResponse(t,recorder)
		})

	}
	
	
}

func TestGetAccount(t *testing.T) {
	account := RandomAccount()

	testcase := []struct {
		name string
		accountID int64
		buildStubs func(store *mockstore.MockStore)
		checkResponse func(t *testing.T,recorder * httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockstore.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(account,nil)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder){
				require.Equal(t,http.StatusOK,recorder.Code)
				requiredBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "Not found",
			accountID: account.ID,
			buildStubs: func(store *mockstore.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(db.Account{},sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder){
				require.Equal(t,http.StatusNotFound,recorder.Code)
			},
		},
		{
			name: "Internal Server",
			accountID: account.ID,
			buildStubs: func(store *mockstore.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(db.Account{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T,recorder * httptest.ResponseRecorder){
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
	}

	for i := range testcase {
		ts := testcase[i]

		t.Run(ts.name,func(t *testing.T) {
		ctrl := gomock.NewController(t)

		store := mockstore.NewMockStore(ctrl)

		ts.buildStubs(store)

		// start test http server
		server := newTestServer(t,store)

		// record info response server of request to server
		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/account/%d",account.ID)

		request,err := http.NewRequest(http.MethodGet,url,nil)

		require.NoError(t,err)

		server.router.ServeHTTP(recorder,request)

		ts.checkResponse(t,recorder)
		})

	}
}

func RandomAccount() db.Account {
	return db.Account{
		ID: utils.RandomInt(1,1000),
		Owner: utils.RandomOwner(),
		Balance: utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}
}

func requiredBodyMatchAccount(t *testing.T,body *bytes.Buffer,account db.Account) {
	data,err := io.ReadAll(body)

	require.NoError(t,err)

	var gotAccount db.Account

	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t,err)

	require.Equal(t,gotAccount,account)
}