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

func (s *pasetoService) CreateToken(payload domain.PasetoPayload, duration time.Duration) (string, error) {
	// Atur waktu terbit dan kedaluwarsa
	now := time.Now()
	exp := now.Add(duration)

	payload.IssuedAt = now
	payload.ExpiresAt = exp

	jsonToken := paseto.JSONToken{
		IssuedAt:   now,
		Expiration: exp,
	}
	// Menambahkan seluruh payload kustom
	jsonToken.Set("uid", payload.UserID)
	jsonToken.Set("eml", payload.Email)
	jsonToken.Set("sid", payload.SessionID)
	jsonToken.Set("rol", payload.RoleName)
	jsonToken.Set("typ", payload.TokenType)
	if payload.VerificationAttemptID != "" {
		jsonToken.Set("vid", payload.VerificationAttemptID)
	}

	// Encrypt (Symmetric)
	return s.paseto.Encrypt(s.symmetricKey, jsonToken, nil)
}

func (s *pasetoService) ValidateToken(tokenString string, expectedType string) (payload domain.PasetoPayload, err error) {
	var jsonToken paseto.JSONToken
	// Inisialisasi payload kosong
	payload = domain.PasetoPayload{}

	// Decrypt (Symmetric)
	err = s.paseto.Decrypt(tokenString, s.symmetricKey, &jsonToken, nil)
	if err != nil {
		err = errors.New("invalid token")
		return
	}

	// Validasi expiration standard
	if err = jsonToken.Validate(); err != nil {
		err = fmt.Errorf("token has expired: %w", err)
		return
	}

	// Ekstrak data dari custom claims
	err = jsonToken.Get("typ", &payload.TokenType)
	if err != nil || payload.TokenType == "" {
		err = errors.New("invalid token: missing token type")
		return
	}
	if payload.TokenType != expectedType {
		err = fmt.Errorf("invalid token type: expected '%s' but got '%s'", expectedType, payload.TokenType)
		return
	}
	err = jsonToken.Get("uid", &payload.UserID)
	if err != nil || payload.UserID == "" {
		err = errors.New("invalid user ID in token")
		return
	}
	err = jsonToken.Get("eml", &payload.Email)
	if err != nil || payload.Email == "" {
		err = errors.New("invalid email in token")
		return
	}
	err = jsonToken.Get("sid", &payload.SessionID)
	if err != nil || payload.SessionID == "" {
		// Handle token lama yang mungkin tidak punya sid, generate baru jika perlu?
		// Atau anggap error jika sid wajib ada. Kita anggap wajib.
		err = errors.New("invalid session ID in token")
		return
	}
	err = jsonToken.Get("rol", &payload.RoleName)
	if err != nil || payload.RoleName == "" {
		err = errors.New("invalid role name in token")
		return
	}
	// Ekstrak attempt ID (opsional)
	_ = jsonToken.Get("vid", &payload.VerificationAttemptID)

	// Isi IssuedAt dan ExpiresAt dari jsonToken
	payload.IssuedAt = jsonToken.IssuedAt
	payload.ExpiresAt = jsonToken.Expiration

	// Validasi format UUID
	_, errUid := uuid.Parse(payload.UserID)
	_, errSid := uuid.Parse(payload.SessionID)
	if errUid != nil || errSid != nil {
		err = errors.New("invalid UUID format in token payload")
		return
	}

	return payload, nil
}
