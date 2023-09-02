package main

import (
	"context"
	"github/wry-0313/exchange/config"
	"github/wry-0313/exchange/db"
	"github/wry-0313/exchange/exchange"
	"github/wry-0313/exchange/user"
	"github/wry-0313/exchange/validator.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main() {

	v := validator.New()

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	
	server := http.Server{
		Addr: cfg.ServerPort,
	}

	userRepo := user.NewRepository(db.DB)

	userService := user.NewService(userRepo, v)

	



	exchangeService := exchange.NewExchange()
	exchangeService.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %s", err)
	}

	close(exchangeService.Shutdown) 
}
