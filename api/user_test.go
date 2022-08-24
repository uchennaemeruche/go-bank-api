package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/uchennaemeruche/go-bank-api/db/mock"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func TestCreateUser(t *testing.T) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := randomUser(hashedPassword)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// build stubs for the mockstore
	body := gin.H{
		"username":  user.Username,
		"password":  password,
		"full_name": user.FullName,
		"email":     user.Email,
	}
	// arg := db.CreateUserParams{
	// 	Username:       user.Username,
	// 	HashedPassword: user.HashedPassword,
	// 	FullName:       user.FullName,
	// 	Email:          user.Email,
	// }

	store.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Times(1).
		Return(user, nil)

	// start a test httpserver
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	data, err := json.Marshal(body)
	require.NoError(t, err)

	request, err := http.NewRequest("POST", "/users", bytes.NewReader(data))
	require.NoError(t, err)

	server.Router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchUser(t, recorder.Body, user)

}

func randomUser(hashedPassword string) db.User {

	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
