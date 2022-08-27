package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type AccountService interface {
	// Validate()
	Create(owner, currency, account_type string) (db.Account, error)
	GetOne(id int64) (db.Account, error)
	ListAccount(owner string, pageSize, pageId int32) ([]db.Account, error)
	UpdateAccount(id, balance int64) (db.Account, error)
	DeleteAccount(id int64) error
}

type service struct {
	store db.Store
}

func NewAccountService(store db.Store) AccountService {
	return &service{
		store: store,
	}
}

func (s *service) Create(owner, currency, account_type string) (db.Account, error) {

	arg := db.CreateAccountParams{
		Owner:       owner,
		Currency:    currency,
		Balance:     0,
		AccountType: account_type,
	}

	account, err := s.store.CreateAccount(context.Background(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				err = fmt.Errorf("with %w", &api.RequestError{Code: 403, Err: err})
				// err = &api.RequestError{
				// 	Code: 403,
				// 	Err:  err,
				// }

			}
		}
	}

	return account, err

}

func (s *service) GetOne(id int64) (db.Account, error) {
	acct, err := s.store.GetAccount(context.Background(), id)

	if err == sql.ErrNoRows {
		err = &api.RequestError{
			Code: 404,
			Err:  errors.New("no record found with the given ID"),
		}
		// err = errors.New("no record found with the given ID")
	}

	return acct, err
}

func (s *service) ListAccount(owner string, pageSize, pageId int32) (accounts []db.Account, err error) {
	arg := db.ListAccountsParams{
		Owner:  owner,
		Limit:  pageSize,
		Offset: (pageId - 1) * pageSize,
	}
	return s.store.ListAccounts(context.Background(), arg)
}

func (s *service) UpdateAccount(id, balance int64) (db.Account, error) {
	arg := db.UpdateAccountParams{
		ID:      id,
		Balance: balance,
	}
	return s.store.UpdateAccount(context.Background(), arg)
}

func (s *service) DeleteAccount(id int64) error {
	return s.store.DeleteAccount(context.Background(), id)

}
