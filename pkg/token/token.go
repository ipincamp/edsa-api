package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto/v2"
)

// ErrInvalidToken digunakan saat token tidak valid
var ErrInvalidToken = errors.New("token is invalid")

// ErrExpiredToken digunakan saat token sudah expired
var ErrExpiredToken = errors.New("token has expired")

// PasetoMaker adalah struct untuk membuat dan verifikasi token Paseto
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// Payload adalah data yang di-encode ke dalam token
type Payload struct {
	TokenID   uuid.UUID `json:"tid"`
	TokenType string    `json:"tty"`
	UserID    string    `json:"uid"`
	RoleID    string    `json:"rid"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"eat"`
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != 32 {
		return nil, errors.New("invalid key size: must be exactly 32 characters")
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func NewPayload(userID string, roleID string, tokenType string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		TokenID:   tokenID,
		TokenType: tokenType,
		UserID:    userID,
		RoleID:    roleID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (maker *PasetoMaker) CreateToken(userID string, roleID string, tokenType string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, roleID, tokenType, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
