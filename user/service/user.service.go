package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserService interface {
	Create(username, hashedPassword, fullname, email string) (db.User, error)
}

type service struct {
	store db.Store
}

func NewUserService(store db.Store) UserService {
	return &service{
		store: store,
	}
}

func (s *service) Create(username, hashedPassword, fullname, email string) (db.User, error) {
	arg := db.CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullname,
		Email:          email,
	}

	user, err := s.store.CreateUser(context.Background(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				var field string
				if pqErr.Constraint == "users_pkey" {
					field = "username"
				} else {
					field = strings.Split(pqErr.Constraint, "_")[1]
					// field = "email"
				}
				err = &api.RequestError{
					Code: 403,
					Err:  fmt.Errorf("provided %s is already in use", field),
				}
			}
		}
	}

	return user, err
}
