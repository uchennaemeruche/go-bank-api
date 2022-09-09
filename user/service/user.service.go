package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/token"
	"github.com/uchennaemeruche/go-bank-api/util"
)

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

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

type RenewAcesTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type RenewAcesTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type DestroySessionReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type ToggleSessionReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	Status       bool   `json:"status" binding:"required"`
}

type LoginResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserResponse
}

type UserService interface {
	Create(username, hashedPassword, fullname, email string) (db.User, error)
	LoginUser(username, password, userAgent, clientIp string, accessTokenDuration, refreshTokenDuration time.Duration) (res LoginResponse, err error)
	GetNewAccessToken(refreshToken string, accessTokenDuration time.Duration) (RenewAcesTokenResponse, error)
	DestroySession(refreshToken string) (bool, error)
	ToggleBlockSession(refreshToken string, status bool) (bool, error)
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

func NewLoginResponse(sessionId uuid.UUID, accessToken, refreshToken string, accessTokenExpiresAt, refreshTokenExpiresAt time.Time, user db.User) LoginResponse {
	return LoginResponse{
		SessionID:             sessionId,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		User:                  NewUserResponse(user),
	}
}

func NewUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
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

func (s *service) LoginUser(username, password, userAgent, clientIp string, accessTokenDuration, refreshTokenDuration time.Duration) (LoginResponse, error) {
	user, err := s.GetUser(username)
	if err != nil {
		return LoginResponse{}, err
	}

	err = util.ComparePassword(password, user.HashedPassword)
	if err != nil {
		err = &api.RequestError{
			Code: 401,
			Err:  errors.New("incorrect login details"),
		}
		return LoginResponse{}, err
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(username, accessTokenDuration)
	if err != nil {
		fmt.Println("Token Error: ", err)
		err = &api.RequestError{
			Code: 500,
			Err:  fmt.Errorf("an internal server error occured while creating user token"),
		}
		return LoginResponse{}, err
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(username, refreshTokenDuration)
	if err != nil {
		fmt.Println("Token Error: ", err)
		err = &api.RequestError{
			Code: 500,
			Err:  fmt.Errorf("an internal server error occured while creating user token"),
		}
		return LoginResponse{}, err
	}

	session, err := s.store.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println("Err Code", pqErr.Code)
			fmt.Println("Err Code Name", pqErr.Code.Name())
			fmt.Println("Err Contraint:", pqErr.Constraint)
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				err = &api.RequestError{
					Code: 500,
					Err:  fmt.Errorf("an internal server error occured while creating user session"),
				}
			}
		}
		return LoginResponse{}, err
	}

	loginRes := NewLoginResponse(
		session.ID, accessToken, refreshToken, accessTokenPayload.ExpiredAt, refreshTokenPayload.ExpiredAt, user,
	)

	return loginRes, nil

}

func (s *service) GetNewAccessToken(refreshToken string, accessTokenDuration time.Duration) (RenewAcesTokenResponse, error) {
	refreshPaylaod, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		err = &api.RequestError{
			Code: 401,
			Err:  err,
		}
		return RenewAcesTokenResponse{}, err
	}

	session, err := s.store.GetSession(context.Background(), refreshPaylaod.ID)
	if err != nil {
		fmt.Println("Session Error: ", err)
		err = &api.RequestError{
			Code: 500,
			Err:  err,
		}
		if err.Error() == sql.ErrNoRows.Error() {
			err = &api.RequestError{
				Code: 404,
				Err:  fmt.Errorf("session does not exist/ invalid session id"),
			}
		}
		return RenewAcesTokenResponse{}, err
	}

	if session.IsBlocked {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("session is blocked, deactivated or inactive"),
		}
		return RenewAcesTokenResponse{}, err
	}

	if session.Username != refreshPaylaod.Username {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("incorrect session user passed"),
		}
		return RenewAcesTokenResponse{}, err
	}

	if session.RefreshToken != refreshToken {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("mismatched session token"),
		}
		return RenewAcesTokenResponse{}, err
	}

	if time.Now().After(session.ExpiresAt) {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("expired session"),
		}
		return RenewAcesTokenResponse{}, err
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(refreshPaylaod.Username, accessTokenDuration)
	if err != nil {
		err = &api.RequestError{
			Code: 500,
			Err:  fmt.Errorf("an internal server error occured while creating user token"),
		}
		return RenewAcesTokenResponse{}, err
	}

	result := RenewAcesTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}
	return result, nil

}

func (s *service) DestroySession(refreshToken string) (bool, error) {
	refreshPayload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		err = &api.RequestError{
			Code: 401,
			Err:  err,
		}
		return false, err
	}

	err = s.store.ExpireSession(context.Background(), db.ExpireSessionParams{
		ID:        refreshPayload.ID,
		ExpiresAt: time.Now(),
	})

	if err != nil {
		err = &api.RequestError{
			Code: 500,
			Err:  fmt.Errorf("could not terminate current session. try again later"),
		}
		return false, err
	}
	return true, nil
}

func (s *service) ToggleBlockSession(refreshToken string, status bool) (bool, error) {

	refreshPayload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		err = &api.RequestError{
			Code: 401,
			Err:  err,
		}
		return false, err
	}
	session, err := s.store.GetSession(context.Background(), refreshPayload.ID)
	if err != nil {
		fmt.Println("Session Error: ", err)
		err = &api.RequestError{
			Code: 500,
			Err:  err,
		}
		if err.Error() == sql.ErrNoRows.Error() {
			err = &api.RequestError{
				Code: 404,
				Err:  fmt.Errorf("session does not exist/ invalid session id"),
			}
		}
		return false, err
	}

	if session.Username != refreshPayload.Username {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("incorrect session user passed"),
		}
		return false, err
	}

	if session.RefreshToken != refreshToken {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("mismatched session token"),
		}
		return false, err
	}

	if time.Now().After(session.ExpiresAt) {
		err = &api.RequestError{
			Code: 401,
			Err:  fmt.Errorf("expired session"),
		}
		return false, err
	}

	err = s.store.ToggleBlockSession(context.Background(), db.ToggleBlockSessionParams{
		ID:        refreshPayload.ID,
		IsBlocked: status,
	})
	if err != nil {
		err = &api.RequestError{
			Code: 500,
			Err:  fmt.Errorf("could not block/unblock user at this time "),
		}
		return false, err
	}
	return true, nil

}
