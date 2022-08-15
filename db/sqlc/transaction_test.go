package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func createRandomTransaction(t *testing.T, account Account) Transaction {
	arg := CreateTransactionParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.AccountID, transaction.AccountID)
	require.Equal(t, arg.Amount, transaction.Amount)

	require.NotZero(t, transaction.AccountID)
	require.NotZero(t, transaction.CreatedAt)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	account := createRandomAccount(t)
	createRandomTransaction(t, account)

}

func TestGetTransaction(t *testing.T) {
	account := createRandomAccount(t)
	transaction := createRandomTransaction(t, account)

	res, err := testQueries.GetTransaction(context.Background(), transaction.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, transaction.ID, res.ID)
	require.Equal(t, transaction.AccountID, res.AccountID)
	require.Equal(t, transaction.Amount, res.Amount)
	require.WithinDuration(t, transaction.CreatedAt, res.CreatedAt, time.Second)
}

func TestListTransaction(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransaction(t, account)
	}

	arg := ListTransactionsParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	transactions, err := testQueries.ListTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
		require.Equal(t, arg.AccountID, transaction.AccountID)
	}
}
