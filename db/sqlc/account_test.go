package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// 常见货币种类
/*
	RMB（人民币）、HKD（港币）、USD（美元）、EUR（欧元）、JPY（日元）、GBP（英镑）
*/

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "xzz",
		Balance:  100,
		Currency: "RMB",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
