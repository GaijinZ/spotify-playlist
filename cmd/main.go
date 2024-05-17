package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"spf-playlist/api/spotify/models"
	"spf-playlist/handler"
	"spf-playlist/pkg/config"
	"spf-playlist/pkg/logger"
	"spf-playlist/pkg/redis"
	"spf-playlist/pkg/sql"
	"spf-playlist/router"
	"spf-playlist/server"
	"spf-playlist/utils"

	spotifyAuth "spf-playlist/api/spotify/auth"
	userAuth "spf-playlist/users/handler/auth"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	var cfg config.GlobalEnv
	var ctx context.Context
	token := &models.Token{}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))
	log := utils.GetLogger(ctx)

	if err := envconfig.Process("spf", &cfg); err != nil {
		log.Fatalf("Failed to process enviromental variables: %v", err)
	}

	DB, err := sql.InitDB(log, cfg, ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	err = redisClient.Ping(ctx)
	if err == nil {
		log.Infof("Redis connected")
	} else {
		log.Infof("Failed to ping to Redis: %v", err)
	}

	newUserAuth := userAuth.NewUserAuth(ctx, cfg, DB, redisClient)
	newSpotifyAuth := spotifyAuth.NewSpotifyAuth(cfg, ctx)
	spotifyHandler := handler.NewSpotifyHandler(*token, ctx, *newSpotifyAuth, cfg)

	r := router.Router(newUserAuth, *spotifyHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Handler: r,
	}

	go server.Run(cfg.Host, cfg.Port, srv, log)

	defer func() {
		cancel()
		redisClient.Close()
		DB.Close()
	}()

	select {
	case <-interrupt:
		fmt.Println("Received a shutdown signal...")
		close(interrupt)
	case <-ctx.Done():
		fmt.Println("Context done")
		close(interrupt)
	}
}
