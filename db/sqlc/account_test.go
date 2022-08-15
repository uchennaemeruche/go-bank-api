package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	res, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, account.ID, res.ID)
	require.Equal(t, account.Owner, res.Owner)
	require.Equal(t, account.Balance, res.Balance)
	require.Equal(t, account.Currency, res.Currency)
	require.WithinDuration(t, account.CreatedAt, res.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	res, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, account.ID, res.ID)
	require.Equal(t, account.ID, res.ID)
	require.Equal(t, account.Owner, res.Owner)
	require.Equal(t, arg.Balance, res.Balance)
	require.Equal(t, account.Currency, res.Currency)
	require.WithinDuration(t, account.CreatedAt, res.CreatedAt, time.Second)

}
