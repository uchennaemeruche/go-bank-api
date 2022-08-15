package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func createRandomTransfer(t *testing.T, sender, recipient Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: sender.ID,
		ToAccountID:   recipient.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer

}

func TestCreateTransfer(t *testing.T) {
	sender := createRandomAccount(t)
	recipient := createRandomAccount(t)

	createRandomTransfer(t, sender, recipient)
}

func TestGetTransfer(t *testing.T) {
	sender := createRandomAccount(t)
	recipient := createRandomAccount(t)
	transfer := createRandomTransfer(t, sender, recipient)

	res, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, transfer.ID, res.ID)
	require.Equal(t, transfer.FromAccountID, res.FromAccountID)
	require.Equal(t, transfer.ToAccountID, res.ToAccountID)
	require.Equal(t, transfer.Amount, res.Amount)
	require.WithinDuration(t, transfer.CreatedAt, res.CreatedAt, time.Second)

}

func TestListTransfers(t *testing.T) {
	sender := createRandomAccount(t)
	recipient := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, sender, recipient)
		createRandomTransfer(t, recipient, sender)
	}

	arg := ListTransfersParams{
		FromAccountID: sender.ID,
		ToAccountID:   sender.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == sender.ID || transfer.ToAccountID == sender.ID)
	}

}
