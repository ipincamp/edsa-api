package util

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

type PasetoPayload struct {
	UserID    string    `json:"uid"`
	RoleID    string    `json:"rid"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(userID string, roleID string, duration time.Duration) (string, error) {
	payload := &PasetoPayload{
		UserID:    userID,
		RoleID:    roleID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*PasetoPayload, error) {
	var payload PasetoPayload
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, fmt.Errorf("token has expired")
	}

	return &payload, nil
}
