package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type AccountService interface {
	// Validate()
	Create(owner, currency string) (db.Account, error)
	GetOne(id int64) (db.Account, error)
	ListAccount(pageSize, pageId int32) ([]db.Account, error)
}

type service struct {
	store *db.Store
}

func NewAccountService(store *db.Store) AccountService {
	return &service{
		store: store,
	}
}

func (s *service) Create(owner, currency string) (db.Account, error) {

	arg := db.CreateAccountParams{
		Owner:    owner,
		Currency: currency,
		Balance:  0,
	}

	return s.store.CreateAccount(context.Background(), arg)

}

func (s *service) GetOne(id int64) (db.Account, error) {
	acct, err := s.store.GetAccount(context.Background(), id)

	if err == sql.ErrNoRows {
		err = &util.RequestError{
			Code: 404,
			Err:  errors.New("no record found with the given ID"),
		}
		// err = errors.New("no record found with the given ID")
	}

	return acct, err
}

func (s *service) ListAccount(pageSize, pageId int32) (accounts []db.Account, err error) {
	arg := db.ListAccountsParams{
		Limit:  pageSize,
		Offset: (pageId - 1) * pageSize,
	}
	return s.store.ListAccounts(context.Background(), arg)
}