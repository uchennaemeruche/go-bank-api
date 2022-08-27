package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/uchennaemeruche/go-bank-api/db/mock"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/token"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs for the mockstore
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)

			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs for the mockstore
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs for the mockstore
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs for the mockstore
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name:      "Unauthorized",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, "unauthorized_username", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs for the mockstore
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// addAuthHeader(t, request, tokenMaker, bearerAuthType, "unauthorized_username", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs for the mockstore
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Start a test httpserver
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func TestCreateAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{"currency": account.Currency, "owner": account.Owner, "account_type": account.AccountType},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:       account.Owner,
					Balance:     0,
					Currency:    account.Currency,
					AccountType: account.AccountType,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{"currency": account.Currency, "owner": account.Owner, "account_type": account.AccountType},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{"currency": "invalid", "owner": account.Owner, "account_type": account.AccountType},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, bearerAuthType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{"currency": account.Currency, "owner": account.Owner, "account_type": account.AccountType},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// addAuthHeader(t, request, tokenMaker, bearerAuthType, "unauthorized", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:       account.Owner,
					Balance:     0,
					Currency:    account.Currency,
					AccountType: account.AccountType,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs for the mockstore
			tc.buildStubs(store)

			// Start a test httpserver
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/accounts", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}

}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:          util.RandomInt(1, 1000),
		Owner:       owner,
		Balance:     util.RandomMoney(),
		Currency:    util.RandomCurrency(),
		AccountType: util.RandomAccountType(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
