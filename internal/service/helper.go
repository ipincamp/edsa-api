package service

import (
	"time"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

func createAuthResponse(user domain.User, s *authService) (dto.AuthResponse, error) {
	tokenTTL := time.Duration(s.config.Paseto.TokenTTLMin) * time.Minute
	pasetoMaker, err := util.NewPasetoMaker(s.config.Paseto.SecretKey)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := pasetoMaker.CreateToken(user.ID, user.Role.ID, tokenTTL)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	res := dto.AuthResponse{
		Token: token,
		User: dto.UserData{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role.Name,
			JoinedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	return res, nil
}
