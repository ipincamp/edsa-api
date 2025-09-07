package dto

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,gte=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	Token Token        `json:"token"`
	User  UserResponse `json:"user"`
}

type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func ToAuthResponse(auth AuthResponse) AuthResponse {
	return AuthResponse{
		Token: Token{
			Access:  auth.Token.Access,
			Refresh: auth.Token.Refresh,
		},
		User: UserResponse{
			ID:        auth.User.ID,
			Name:      auth.User.Name,
			Email:     auth.User.Email,
			Role:      auth.User.Role,
			JoinedAt:  auth.User.JoinedAt,
			UpdatedAt: auth.User.UpdatedAt,
		},
	}
}
