package utils

import (
	"errors"
	"time"

	"github.com/o1egl/paseto/v2"

	"github.com/ipincamp/edsa/internal/models"
)

type PasetoPayload struct {
	UserID    string    `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

func NewPayload(user *models.User, duration time.Duration) (*PasetoPayload, error) {
	payload := &PasetoPayload{
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func CreateToken(payload *PasetoPayload, symmetricKey string) (string, error) {
	return paseto.NewV2().Encrypt([]byte(symmetricKey), payload, nil)
}

func VerifyToken(token string, symmetricKey string) (*PasetoPayload, error) {
	payload := &PasetoPayload{}
	err := paseto.NewV2().Decrypt(token, []byte(symmetricKey), payload, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, errors.New("token has expired")
	}

	return payload, nil
}
