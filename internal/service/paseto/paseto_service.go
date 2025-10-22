package paseto

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"github.com/o1egl/paseto/v2"
)

type pasetoService struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

// Kustom Paseto payload
type PasetoPayload struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
}

func NewPasetoService(symmetricKeyBase64 string) (usecase.TokenService, error) {
	if len(symmetricKeyBase64) == 0 {
		return nil, errors.New("PASETO_SYMMETRIC_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(symmetricKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode paseto key: %w", err)
	}

	if len(key) != 32 {
		return nil, errors.New("paseto key must be 32 bytes")
	}

	return &pasetoService{
		symmetricKey: key,
		paseto:       paseto.NewV2(),
	}, nil
}

func (s *pasetoService) CreateToken(user *domain.User, sessionID uuid.UUID, duration time.Duration) (string, error) {
	payload := PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: sessionID.String(),
	}

	// Token berlaku selama 24 jam
	now := time.Now()
	exp := now.Add(duration)

	jsonToken := paseto.JSONToken{
		IssuedAt:   now,
		Expiration: exp,
	}
	// Menambahkan payload kustom
	jsonToken.Set("data", payload)

	// Encrypt (Symmetric)
	return s.paseto.Encrypt(s.symmetricKey, jsonToken, nil)
}

func (s *pasetoService) ValidateToken(tokenString string) (uuid.UUID, uuid.UUID, error) {
	var jsonToken paseto.JSONToken
	var payload PasetoPayload

	// Decrypt (Symmetric)
	err := s.paseto.Decrypt(tokenString, s.symmetricKey, &jsonToken, nil)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid token")
	}

	// Validasi expiration
	if err := jsonToken.Validate(); err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("token has expired: %w", err)
	}

	// Ekstrak payload kustom
	if err := jsonToken.Get("data", &payload); err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get payload from token: %w", err)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user ID format in token payload")
	}

	sessionID, err := uuid.Parse(payload.SessionID)
	if err != nil {
		// Jika token lama (sebelum update ini) tidak memiliki sessionID, buatkan yang baru
		sessionID = uuid.New()
	}

	return userID, sessionID, nil
}
