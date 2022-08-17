package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to run DB queries individually and withing a transaction
type Store struct {
	// Composition: a preferred way of extending struct functionality. By embedding Queries in Store struct, all individual query functions provided by Queries will be available to Store.
	// Moreso, we can define new functions for Store struct
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer        Transfer    `json:"transfer"`
	FromAccount     Account     `json:"from_account"`
	ToAccount       Account     `json:"to_account"`
	FromTransaction Transaction `json:"from_transaction"`
	ToTransaction   Transaction `json:"to_transaction"`
}

var txKey = struct{}{}

// TransferTx performs money transfer from one account to another using the below steps.
// 1. Create a Transfer record
// 2. Add Account Transactions
// 3. Update Account balance
// All within a single transaction

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		// 1. Create a Transfer record
		fmt.Println(txName, "Create Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. Add Account Transaction entry 1
		fmt.Println(txName, "Create Transaction 1")
		result.FromTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. Add Account Transaction counter entry 2
		fmt.Println(txName, "Create Transaction 2")
		result.ToTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update Account Balance
		// fmt.Println(txName, "Get Sender Account ")
		// senderAcct, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "Update Sender Account ")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: senderAcct.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Get Recipient Account ")
		// recipientAcct, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "Update Account ")
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountId,
		// 	Balance: recipientAcct.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}
