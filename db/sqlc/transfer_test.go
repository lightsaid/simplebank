package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"lightsaid.com/simplebank/utils"
)

func createRandomTransfer(t *testing.T, a1, a2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        utils.RandomMoney(),
	}
	tfr, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tfr)

	require.Equal(t, tfr.FromAccountID, a1.ID)
	require.Equal(t, tfr.ToAccountID, a2.ID)
	require.Equal(t, tfr.Amount, arg.Amount)

	return tfr
}

func TestCreateTransfer(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	_ = createRandomTransfer(t, a1, a2)
}

func TestGetTransfer(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	tf1 := createRandomTransfer(t, a1, a2)

	tf2, err := testQueries.GetTransfer(context.Background(), tf1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tf2)

	EqualStruct(t, tf1, tf2)
}

func TestListTransfer(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	// 相互转账5次
	for i := 0; i < 5; i++ {
		_ = createRandomTransfer(t, a1, a2)
		_ = createRandomTransfer(t, a2, a1)
	}

	// 获取列表
	arg := ListTransfersParams{
		FromAccountID: a1.ID,
		ToAccountID:   a1.ID,
		Limit:         5,
		Offset:        5,
	}

	tts, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tts, 5)

	for _, tt := range tts {
		require.NotEmpty(t, tt)
		require.True(t, tt.FromAccountID == a1.ID || tt.ToAccountID == a1.ID)
	}
}
