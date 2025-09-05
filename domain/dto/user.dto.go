package dto

type UpdateProfileUserRequest struct {
	Name                    string `json:"name,omitempty" validate:"omitempty,min=3"`
	OldPassword             string `json:"old_password,omitempty" validate:"omitempty"`
	NewPassword             string `json:"new_password,omitempty" validate:"omitempty,gte=8"`
	NewPasswordConfirmation string `json:"new_password_confirmation,omitempty" validate:"omitempty,eqfield=NewPassword"`
}

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	JoinedAt  string `json:"joined_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
