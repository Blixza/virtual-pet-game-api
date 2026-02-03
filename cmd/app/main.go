package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"virtual_pet_game/config"
	db_config "virtual_pet_game/config/db"
	domain_pet "virtual_pet_game/internal/domain/pet"
	domain_user "virtual_pet_game/internal/domain/user"
	auth_handler "virtual_pet_game/internal/handler/auth"
	handler_pet "virtual_pet_game/internal/handler/pet"
	"virtual_pet_game/internal/middleware"
	service_pet "virtual_pet_game/internal/service/pet"
	"virtual_pet_game/pkg/db"
	"virtual_pet_game/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	envPath := flag.String("envpath", ".env", "Path to .env file")
	dbEnvPath := flag.String("dbenvpath", ".env", "Path to db .env file")

	flag.Parse()

	cfg := config.Load(*envPath)

	log := logger.New(cfg.Level)

	ctx := context.Background()

	dbCfg := db_config.Load(*dbEnvPath)

	dsn := db.NewDSN(&dbCfg)
	db := db.NewPool(ctx, dsn)
	defer db.Close()

	petRepo := domain_pet.NewRepository(db, log)
	petService := service_pet.New(petRepo)
	petHandler := handler_pet.New(petService, log)

	userRepo := domain_user.NewRepository(db)
	authHandler := auth_handler.New(userRepo, cfg.Secret, log)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("pong")
	})

	mux.Handle("POST /pets", middleware.AuthMiddleware(http.HandlerFunc(petHandler.Create), cfg.Secret))
	mux.Handle("GET /pets/id/{id}", middleware.AuthMiddleware(http.HandlerFunc(petHandler.Get), cfg.Secret))
	mux.Handle("GET /pets/name/{name}", middleware.AuthMiddleware(http.HandlerFunc(petHandler.Get), cfg.Secret))
	mux.Handle("PUT /pets/id/{id}", middleware.AuthMiddleware(http.HandlerFunc(petHandler.Update), cfg.Secret))
	mux.Handle("DELETE /pets/id/{id}", middleware.AuthMiddleware(http.HandlerFunc(petHandler.Delete), cfg.Secret))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HttpPort),
		Handler: mux,
	}

	go func() {
		log.Info("server starting", zap.String("addr", srv.Addr))
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal("server start failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited properly")
}
