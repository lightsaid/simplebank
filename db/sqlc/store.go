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

		// NOTE: 更新 balance 方式 1， 执行2条sql不高效率 (同时并发 select * from ... for update 会死锁)
		// // 4. a 账户 balance - 10
		// account1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// updateAccount1, err := q.UpdateAccount(context.Background(), UpdateAccountParams{
		// 	ID:      account1.ID,
		// 	Balance: account1.Balance - arg.Amount,
		// })

		// if err != nil {
		// 	return err
		// }

		// // 5. b 账户 balance + 10
		// account2, err := q.GetAccountForUpdate(context.Background(), arg.ToAccountID)

		// if err != nil {
		// 	return err
		// }
		// updateAccount2, err := q.UpdateAccount(context.Background(), UpdateAccountParams{
		// 	ID:      account2.ID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// result.FromAccount = updateAccount1
		// result.ToAccount = updateAccount2

		// NOTE: 更新 balance 方式 2, 一条SQL语句， 没用 for update（排它锁）没用锁，就不会有死锁。
		// if result.FromAccount, err = updateMoney(ctx, q, arg.FromAccountID, -arg.Amount); err != nil {
		// 	return err
		// }
		// if result.ToAccount, err = updateMoney(ctx, q, arg.ToAccountID, arg.Amount); err != nil {
		// 	return err
		// }

		// NOTE: 更新 balance 同时解决双向转账并发问题，有序执行
		if arg.FromAccountID < arg.ToAccountID {
			if result.FromAccount, err = updateMoney(ctx, q, arg.FromAccountID, -arg.Amount); err != nil {
				return err
			}
			if result.ToAccount, err = updateMoney(ctx, q, arg.ToAccountID, arg.Amount); err != nil {
				return err
			}
		} else {
			if result.ToAccount, err = updateMoney(ctx, q, arg.ToAccountID, arg.Amount); err != nil {
				return err
			}
			if result.FromAccount, err = updateMoney(ctx, q, arg.FromAccountID, -arg.Amount); err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func updateMoney(ctx context.Context, q *Queries, accountID, amount int64) (Account, error) {
	param := AddAccountBalanceParams{
		Amount: amount,
		ID:     accountID,
	}
	return q.AddAccountBalance(context.Background(), param)
}
