package dto

type UpdateUserRequest struct {
	Name                    string `json:"name,omitempty" validate:"omitempty,min=3"`
	OldPassword             string `json:"old_password,omitempty" validate:"omitempty"`
	NewPassword             string `json:"new_password,omitempty" validate:"omitempty,gte=8"`
	NewPasswordConfirmation string `json:"new_password_confirmation,omitempty" validate:"omitempty,eqfield=NewPassword"`
}

type UserData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
	UpdatedAt string `json:"updated_at"`
}
