package dto

type UserRequest struct {
	Name     string `json:"name"     validate:"required"`
	Surname  string `json:"surname"  validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role"     validate:"required,oneof=admin lab_assistant lab_manager repairman"`
}

type PatchUserRequest struct {
	Name     *string `json:"name"     validate:"omitempty,min=1"`
	Surname  *string `json:"surname"  validate:"omitempty,min=1"`
	Email    *string `json:"email"    validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	Role     *string `json:"role"     validate:"omitempty,oneof=lab_assistant lab_manager repairman"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
