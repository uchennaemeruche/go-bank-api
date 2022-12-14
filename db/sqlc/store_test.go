package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	sender := createRandomAccount(t)
	recipient := createRandomAccount(t)

	// Run n concurrent transfer transactions
	n := 6
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: sender.ID,
				ToAccountId:   recipient.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Write Tests for Results
	existed := make(map[int]bool)
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

		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// Check the final update balances of two accounts.
	updatedSenderAcct, err := testQueries.GetAccount(context.Background(), sender.ID)
	require.NoError(t, err)

	updatedRecipientAcct, err := testQueries.GetAccount(context.Background(), recipient.ID)
	require.NoError(t, err)

	require.Equal(t, sender.Balance-int64(n)*amount, updatedSenderAcct.Balance)
	require.Equal(t, recipient.Balance+int64(n)*amount, updatedRecipientAcct.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		toAccountID := account1.ID
		fromAccountID := account2.ID

		if i%2 == 1 {
			toAccountID = account2.ID
			fromAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})
			errs <- err

		}()
	}

	// Write Tests for Results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final update balances of two accounts.
	updatedSenderAcct, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedRecipientAcct, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("after >>>", updatedSenderAcct.Balance, updatedRecipientAcct.Balance)

	require.Equal(t, account1.Balance, updatedSenderAcct.Balance)
	require.Equal(t, account2.Balance, updatedRecipientAcct.Balance)

}
