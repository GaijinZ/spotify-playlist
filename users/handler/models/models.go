package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty" form:"email" validate:"required,email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	IsActive bool   `json:"is_active,omitempty"`
}

type UserResponse struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty" validate:"required,email"`
	Role  string `json:"role,omitempty"`
}

type Authentication struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
