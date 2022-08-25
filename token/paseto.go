package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// implements Maker Interface
type PasetoMaker struct {
	paseto      *paseto.V2
	symmeticKey []byte
}

func NewPasetoMaker(symmeticKey string) (Maker, error) {
	if len(symmeticKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be atleat %d characters", secretKeyMinLength)
	}
	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		symmeticKey: []byte(symmeticKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (token string, err error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmeticKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmeticKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
