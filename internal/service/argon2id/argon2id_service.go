package argon2id

import (
	"github.com/alexedwards/argon2id"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type passwordService struct{}

func NewPasswordService() usecase.PasswordService {
	return &passwordService{}
}

func (s *passwordService) Hash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *passwordService) Compare(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
