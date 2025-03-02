package db

import (
	"context"
	"testing"
	"time"

	"github.com/ernesto-rod/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	randAccount := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: randAccount.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	createdEntry := createRandomEntry(t)
	retrievedEntry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntry)

	require.Equal(t, createdEntry.ID, retrievedEntry.ID)
	require.Equal(t, createdEntry.AccountID, retrievedEntry.AccountID)
	require.Equal(t, createdEntry.Amount, retrievedEntry.Amount)
	require.WithinDuration(t, createdEntry.CratedAt, retrievedEntry.CratedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	for range 10 {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
