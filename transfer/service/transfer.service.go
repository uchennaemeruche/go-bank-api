package service

import (
	"context"
	"database/sql"
	"fmt"

	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type TransferService interface {
	Create(fromAccountId, toAccountId, amount int64, currency string) (db.TransferTxResult, error)
	ValidAccount(accountId int64, currency string) (account db.Account, err error)
}

type service struct {
	store db.Store
}

func NewTransferService(store db.Store) TransferService {
	return &service{
		store: store,
	}
}

func (s *service) Create(fromAccountId, toAccountId, amount int64, currency string) (db.TransferTxResult, error) {
	arg := db.TransferTxParams{
		FromAccountID: fromAccountId,
		ToAccountId:   toAccountId,
		Amount:        amount,
	}
	return s.store.TransferTx(context.Background(), arg)
}

func (s *service) ValidAccount(accountId int64, currency string) (account db.Account, err error) {
	account, err = s.store.GetAccount(context.Background(), accountId)

	if err == sql.ErrNoRows {
		err = &api.RequestError{
			Code: 404,
			Err:  fmt.Errorf("invalid account number: %d ", accountId),
		}
	}

	return account, err
}
