package main

import (
	"context"
	"github/wry-0313/exchange/auth"
	"github/wry-0313/exchange/config"
	"github/wry-0313/exchange/db"
	"github/wry-0313/exchange/endpoint"
	"github/wry-0313/exchange/exchange"
	"github/wry-0313/exchange/jwt"
	"github/wry-0313/exchange/middleware"
	"github/wry-0313/exchange/user"
	"github/wry-0313/exchange/validator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)


func main() {

	validator := validator.New()

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	
	// Setup server
	r := chi.NewRouter()
	server := http.Server{
		Addr:    cfg.ServerPort,
		Handler: setupHandler(r, db, validator, cfg),
	}

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

// setupHandler sets up all the middleware and API routes for the server.
func setupHandler(
	r chi.Router,
	db *db.DB,
	v validator.Validate,
	cfg *config.Config,
) chi.Router {
	// Set up middleware
	r.Use(middleware.Cors())

	// Set up repositories
	userRepo := user.NewRepository(db.DB)

	// Set up services
	jwtService := jwt.NewService(cfg.JwtSecret, cfg.JwtExpiration)
	authService := auth.NewService(userRepo, jwtService, v)
	userService := user.NewService(userRepo, v)

	// rdb := ws.NewRedis(cfg.Rdb)

	// Set up API
	userAPI := user.NewAPI(userService, jwtService, v)
	authAPI := auth.NewAPI(authService, v)


	// Set up auth handler
	authHandler := middleware.Auth(jwtService)

	// Register handlers
	userAPI.RegisterHandlers(r, authHandler)
	authAPI.RegisterHandlers(r)

	r.Get("/ping", handlePingCheck)

	return r
}

func handlePingCheck(w http.ResponseWriter, _ *http.Request) {
	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "pong"})
}