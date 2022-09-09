package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uchennaemeruche/go-bank-api/token"
)

func createRandomSession(t *testing.T) Session {

	user := createRandomUser(t)

	uid, err := uuid.NewRandom()
	require.NoError(t, err)

	tokenMaker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	refreshToken, refreshPayload, err := tokenMaker.CreateToken(user.Username, time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, refreshPayload)

	arg := CreateSessionParams{
		ID:           uid,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)

	require.False(t, session.ExpiresAt.IsZero())
	require.NotZero(t, session.CreatedAt)

	return session
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session := createRandomSession(t)
	result, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, session.ID, result.ID)
	require.Equal(t, session.RefreshToken, result.RefreshToken)
	require.Equal(t, session.Username, result.Username)
	require.Equal(t, session.UserAgent, result.UserAgent)
	require.Equal(t, session.ClientIp, result.ClientIp)
	require.Equal(t, session.IsBlocked, result.IsBlocked)

	require.WithinDuration(t, session.CreatedAt, result.CreatedAt, time.Second)
}
