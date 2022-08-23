package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:       user.Username,
		Balance:     util.RandomMoney(),
		Currency:    util.RandomCurrency(),
		AccountType: util.RandomAccountType(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.AccountType, account.AccountType)

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
	require.Equal(t, account.AccountType, res.AccountType)
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
	require.Equal(t, account.AccountType, res.AccountType)
	require.WithinDuration(t, account.CreatedAt, res.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
