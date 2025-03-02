package db

import (
	"context"
	"testing"
	"time"

	"github.com/ernesto-rod/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	createdTransfer := createRandomTransfer(t)
	retrievedTransfer, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, retrievedTransfer)

	require.Equal(t, createdTransfer.ID, retrievedTransfer.ID)
	require.Equal(t, createdTransfer.FromAccountID, retrievedTransfer.FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, retrievedTransfer.ToAccountID)
	require.Equal(t, createdTransfer.Amount, retrievedTransfer.Amount)
	require.WithinDuration(t, createdTransfer.CreatedAt, retrievedTransfer.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for range 10 {
		arg := CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        util.RandomMoney(),
		}

		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
