package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"virtual_pet_game/config"
	db_config "virtual_pet_game/config/db"
	domain_pet "virtual_pet_game/internal/domain/pet"
	handler_pet "virtual_pet_game/internal/handler/pet"
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

	mux := http.NewServeMux()
	mux.HandleFunc("POST /pets", petHandler.Create)
	mux.HandleFunc("GET /pets/id/{id}", petHandler.Get)
	mux.HandleFunc("GET /pets/name/{name}", petHandler.Get)
	mux.HandleFunc("PUT /pets/id/{id}", petHandler.Update)

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
