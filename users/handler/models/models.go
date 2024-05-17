package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
)

type User struct {
	ID       gocql.UUID `json:"id,omitempty"`
	Name     string     `json:"name,omitempty"`
	Email    string     `json:"email,omitempty" form:"email" validate:"required,email"`
	Password string     `json:"password,omitempty"`
	Role     string     `json:"role,omitempty"`
	IsActive bool       `json:"is_active,omitempty"`
}

type UserResponse struct {
	ID    gocql.UUID `json:"id,omitempty"`
	Name  string     `json:"name,omitempty"`
	Email string     `json:"email,omitempty" validate:"required,email"`
	Role  string     `json:"role,omitempty"`
}

type Authentication struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type Claims struct {
	UserID      gocql.UUID `json:"id"`
	Email       string     `json:"email"`
	TokenString string     `json:"token"`
	Role        string     `json:"role"`
	jwt.StandardClaims
}
