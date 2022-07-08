package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"lightsaid.com/simplebank/utils"
)

func createRandomEntry(t *testing.T) Entry {
	a1 := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: a1.ID,
		Amount:    utils.RandomMoney(),
	}
	e1, err := testQueries.CreateEntry(context.Background(), arg)
	require.NotEmpty(t, a1)
	require.NoError(t, err)
	require.Equal(t, e1.AccountID, a1.ID)
	require.Equal(t, e1.Amount, arg.Amount)

	return e1

}

func TestCreateEntry(t *testing.T) {
	_ = createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	e1 := createRandomEntry(t)
	e2, err := testQueries.GetEntry(context.Background(), e1.ID)
	require.NoError(t, err)
	EqualStruct(t, e1, e2)
}

func TestListentries(t *testing.T) {
	a1 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			AccountID: a1.ID,
			Amount:    utils.RandomMoney(),
		}
		testQueries.CreateEntry(context.Background(), arg)
	}

	arg := GetEntriesParams{
		AccountID: a1.ID,
		Limit:     5,
		Offset:    5,
	}

	es, err := testQueries.GetEntries(context.Background(), arg)
	require.NoError(t, err)
	for _, e := range es {
		require.NotEmpty(t, e)
	}
}
