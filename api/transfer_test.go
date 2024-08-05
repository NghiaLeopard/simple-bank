package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockstore "github.com/NghiaLeopard/simple-bank/db/mock"
	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1 := RandomAccount()
	account2 := RandomAccount()

	amount := int64(10)

	arg := db.TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: amount,
	}

	data := struct{
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID int64 `json:"to_account_id"`
		Amount int64 `json:"amount"`
		Currency string `json:"currency"`
	}{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: account1.Balance,
		Currency: "USD",
	}

	ctl := gomock.NewController(t)

	store := mockstore.NewMockStore(ctl)

	store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account1.ID)).Times(1).Return(account1,nil)
	store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account2.ID)).Times(1).Return(account2,nil)

	store.EXPECT().TransferTx(gomock.Any(),gomock.Eq(arg)).Times(1)

	server := newTestServer(t,store)

	recorder := httptest.NewRecorder()

	dataParse,err := utils.JsonReaderFactory(data)

	if err != nil {
		require.NoError(t,err)
	}

	request,err := http.NewRequest(http.MethodPost,"/transfer",dataParse)

	require.NoError(t,err)

	server.router.ServeHTTP(recorder,request)

	require.Equal(t,http.StatusOK,recorder.Code)
}