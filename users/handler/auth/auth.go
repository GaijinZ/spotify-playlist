package auth

import (
	"context"
	"encoding/json"
	"github.com/gocql/gocql"
	"net/http"
	"time"

	"spf-playlist/pkg/config"
	"spf-playlist/pkg/hash"
	"spf-playlist/pkg/middleware"
	"spf-playlist/pkg/redis"
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
	ctx   context.Context
	cfg   config.GlobalEnv
	DB    sql.DBer
	redis *redis.Client
}

func NewUserAuth(ctx context.Context, cfg config.GlobalEnv, DB sql.DBer, redis *redis.Client) *UserAuth {
	return &UserAuth{
		ctx:   ctx,
		cfg:   cfg,
		DB:    DB,
		redis: redis,
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

	err = u.DB.Insert(&user)
	if err != nil {
		log.Errorf("Error inserting user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *UserAuth) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var userAuth models.Authentication

	log := utils.GetLogger(u.ctx)

	err := json.NewDecoder(r.Body).Decode(&userAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	query := u.DB.Get(userAuth.Email)
	err = query.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Errorf("Error getting user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.ctx = context.WithValue(u.ctx, "userID", user.ID)

	comparePassword := hash.ComparePasswords(user.Password, userAuth.Password)
	if !comparePassword {
		log.Errorf("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token, err := middleware.GenerateJWT(user, u.cfg)
	if err != nil {
		log.Errorf("Error generating token: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.redis.Client.Set(u.ctx, user.ID.String(), token, time.Hour*24).Err()
	if err != nil {
		log.Errorf("Error setting token: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof("User has been logged in: %v", user.Email)
	w.WriteHeader(http.StatusOK)
}

func (u *UserAuth) Logout(w http.ResponseWriter, r *http.Request) {
	log := utils.GetLogger(u.ctx)

	userID, ok := u.ctx.Value("userID").(gocql.UUID)
	if !ok {
		http.Error(w, "Invalid userID in context", http.StatusInternalServerError)
		return
	}

	userIDStr := userID.String()

	if err := u.redis.Client.Del(u.ctx, userIDStr).Err(); err != nil {
		log.Errorf("Error deleting token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("User with userID %s has been logged out", userIDStr)
	w.WriteHeader(http.StatusOK)
}
