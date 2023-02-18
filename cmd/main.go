package main

import (
	"context"
	"github.com/Hymiside/evraz-api/kafka_consume"
	"github.com/Hymiside/evraz-api/pkg/handler"
	"github.com/Hymiside/evraz-api/pkg/models"
	"github.com/Hymiside/evraz-api/pkg/repository"
	"github.com/Hymiside/evraz-api/pkg/server"
	"github.com/Hymiside/evraz-api/pkg/service"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Panic("error .env file not found")
	}

	db, err := repository.NewPostgresDB(
		ctx,
		models.ConfigRepository{
			Host:     os.Getenv("HOST_DB"),
			Port:     os.Getenv("PORT_DB"),
			User:     os.Getenv("USER_DB"),
			Password: os.Getenv("PASSWORD_DB"),
			Name:     os.Getenv("NAME_DB"),
		})
	if err != nil {
		log.Panicf("error to init repository: %v", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	go kafka_consume.InitConsume(services)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	srv := server.Server{}
	if err = srv.RunServer(ctx, handlers.InitRoutes(), models.ConfigServer{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT")}); err != nil {

		log.Panicf("failed to run server: %v", err)
	}
}
