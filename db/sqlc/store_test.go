package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// NOTE: 测试驱动 TDD

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	fmt.Println(">> before: ", a1.Balance)

	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		// 检查 transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, a1.ID, transfer.FromAccountID)
		require.Equal(t, a2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// 检查 entry 表
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, a1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, a2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// 检查 account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, a1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, a2.ID, toAccount.ID)

		diff1 := a1.Balance - fromAccount.Balance
		// NOTE: 测试不通过，因为转账存在并发，还没有解决
		fmt.Printf("diff1: %d - %d = %d\n", a1.Balance, fromAccount.Balance, diff1)
		diff2 := toAccount.Balance - a2.Balance
		fmt.Println(">> after: ", fromAccount.Balance)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // diff1 = 1*amount, 2*amount, 3*amount ...

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)

	updateAccount2, err := testQueries.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)

	require.Equal(t, a1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, a2.Balance+int64(n)*amount, updateAccount2.Balance)

}
