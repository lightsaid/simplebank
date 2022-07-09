package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store 提供一个基本的查询和事务功能结构体
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore 创建一个Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx 执行数据库事务操作
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 开启事务
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// 创建一个新的 CRUD 的 Queries
	// sql.DB 和 sql.Tx 都实现了 DBTX interface
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams 转帐事务入参
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `josn:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx 转帐事务
/**
NOTE: 转账交易，假设 a账户 往 b账户 转 10元
1. 创建一条交易记录(transfer)，amount = 10
2. a 账户创建一条转帐记录(entry) = -10
3. b 账户创建一条转帐记录（entry） = +10
4. a 账户 balance - 10
5. b 账户 balance + 10
*/
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err error
	err = store.execTx(ctx, func(q *Queries) error {
		// 实现上面转帐 5 步
		// 1. 创建一条交易记录(transfer)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. a 账户创建一条转帐记录(entry) = -10
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3. b 账户创建一条转帐记录（entry） = +10
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO:
		// 4. a 账户 balance - 10

		// 5. b 账户 balance + 10

		return nil
	})

	return result, err
}
