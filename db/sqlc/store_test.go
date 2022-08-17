package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	sender := createRandomAccount(t)
	recipient := createRandomAccount(t)

	// Run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: sender.ID,
				ToAccountId:   recipient.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Write Tests for Results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check Transfer status
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, sender.ID, transfer.FromAccountID)
		require.Equal(t, recipient.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check
		fromTnx := result.FromTransaction
		require.NotEmpty(t, fromTnx)
		require.Equal(t, sender.ID, fromTnx.AccountID)
		require.Equal(t, -amount, fromTnx.Amount)
		require.NotZero(t, fromTnx.ID)
		require.NotZero(t, fromTnx.CreatedAt)

		_, err = store.GetTransaction(context.Background(), fromTnx.ID)
		require.NoError(t, err)

		toTnx := result.ToTransaction
		require.NotEmpty(t, toTnx)
		require.Equal(t, recipient.ID, toTnx.AccountID)
		require.Equal(t, amount, toTnx.Amount)
		require.NotZero(t, toTnx.ID)
		require.NotZero(t, toTnx.CreatedAt)

		_, err = store.GetTransaction(context.Background(), toTnx.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, sender.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, recipient.ID, toAccount.ID)

		// TODO: Check account Balance
		diff1 := sender.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - recipient.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

	}

}
