package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"spf-playlist/pkg/hash"
	"spf-playlist/pkg/sql"
	"spf-playlist/users/handler/models"
	"spf-playlist/utils"
)

type UserAuther interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type UserAuth struct {
	ctx context.Context
	DB  sql.DBer
}

func NewUserAuth(DB sql.DBer, ctx context.Context) *UserAuth {
	return &UserAuth{
		ctx: ctx,
		DB:  DB,
	}
}

func (u *UserAuth) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	log := utils.GetLogger(u.ctx)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user.Password, err = hash.GenerateHashPassword(user.Password)
	if err != nil {
		log.Errorf("Error hashing password: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_ = u.DB.Insert(&user)

}

func (u *UserAuth) Login(w http.ResponseWriter, r *http.Request) {}

func (u *UserAuth) Logout(w http.ResponseWriter, r *http.Request) {}
