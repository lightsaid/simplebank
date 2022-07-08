package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"lightsaid.com/simplebank/utils"
)

// 创建随机账号
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	a1 := createRandomAccount(t)
	a2, err := testQueries.GetAccount(context.Background(), a1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, a2)

	require.Equal(t, a1.ID, a2.ID)
	require.Equal(t, a1.Balance, a2.Balance)
	require.Equal(t, a1.Currency, a2.Currency)
	require.Equal(t, a1.Owner, a2.Owner)
	require.WithinDuration(t, a1.CreatedAt, a2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	a1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      a1.ID,
		Balance: utils.RandomMoney(),
	}

	a2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, a1.ID, a2.ID)
	require.Equal(t, arg.Balance, a2.Balance)
	require.Equal(t, a1.Owner, a2.Owner)
	require.Equal(t, a1.Currency, a2.Currency)
}

func TestDeleteAccount(t *testing.T) {
	a1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), a1.ID)

	require.NoError(t, err)

	a2, err := testQueries.GetAccount(context.Background(), a1.ID)
	require.Empty(t, a2)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
