package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/token"
	"github.com/uchennaemeruche/go-bank-api/util"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
type UserLoginReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserService interface {
	Create(username, hashedPassword, fullname, email string) (db.User, error)
	LoginUser(username, password string, duration time.Duration) (accessToken string, user db.User, err error)
}

type service struct {
	store      db.Store
	tokenMaker token.Maker
}

func NewUserService(store db.Store, tokenMaker token.Maker) UserService {
	return &service{
		store:      store,
		tokenMaker: tokenMaker,
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

func (s *service) GetUser(username string) (db.User, error) {
	user, err := s.store.GetUser(context.Background(), username)

	if err == sql.ErrNoRows {
		err = &api.RequestError{
			Code: 404,
			Err:  errors.New("no record found with the given ID"),
		}
		// err = errors.New("no record found with the given ID")
	}

	return user, err
}

func (s *service) LoginUser(username, password string, duration time.Duration) (accessToken string, user db.User, err error) {
	user, err = s.GetUser(username)
	if err != nil {
		return "", db.User{}, err
	}

	err = util.ComparePassword(password, user.HashedPassword)
	if err != nil {
		err = &api.RequestError{
			Code: 401,
			Err:  errors.New("incorrect login details"),
		}
		return "", db.User{}, err
	}

	acccessToken, err := s.tokenMaker.CreateToken(username, duration)
	if err != nil {
		return "", db.User{}, err
	}

	return acccessToken, user, err

}
