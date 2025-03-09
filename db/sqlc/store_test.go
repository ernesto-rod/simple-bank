package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	srcAccount := createRandomAccount(t)
	dstAccount := createRandomAccount(t)
	fmt.Println(">> before:", srcAccount.Balance, dstAccount.Balance)
	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for range n {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: srcAccount.ID,
				ToAccountID:   dstAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for range n {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, srcAccount.ID, transfer.FromAccountID)
		require.Equal(t, dstAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, srcAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, dstAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, srcAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, dstAccount.ID, toAccount.ID)

		//check account's balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := srcAccount.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - dstAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), srcAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), dstAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", srcAccount.Balance, dstAccount.Balance)
	require.Equal(t, srcAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, dstAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	srcAccount := createRandomAccount(t)
	dstAccount := createRandomAccount(t)
	fmt.Println(">> before:", srcAccount.Balance, dstAccount.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := range n {
		fromAccountID := srcAccount.ID
		toAccountID := dstAccount.ID

		if i%2 == 1 {
			fromAccountID = dstAccount.ID
			toAccountID = srcAccount.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check results
	for range n {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), srcAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), dstAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", srcAccount.Balance, dstAccount.Balance)
	require.Equal(t, srcAccount.Balance, updatedFromAccount.Balance)
	require.Equal(t, dstAccount.Balance, updatedToAccount.Balance)
}
