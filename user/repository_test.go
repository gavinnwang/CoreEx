package user

import (
	"fmt"
	"github/wry-0313/exchange/config"
	"github/wry-0313/exchange/db"
	// "github/wry-0313/exchange/models"
	"testing"
	// "time"
)

func TestRepository(t *testing.T) {

	cfg, err := config.Load(".env")
	if err != nil {
		t.Fatalf("Could not load config: %v", err)
	}

	db, err := db.New(cfg.DB)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := NewRepository(db.DB)

	email := "testing@gmail.com"
	// password := "test_password"
	// err = userRepo.CreateUser(models.User{
	// 	Name:      "Test User",
	// 	Email:     &email,
	// 	Password:  &password,
	// 	UpdatedAt: time.Now(),
	// 	CreatedAt: time.Now(),
	// })

	// if err != nil {
	// 	t.Errorf("Failed to create user: %v", err)
	// }

	user,err := userRepo.GetUserByEmail(email)

	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	fmt.Printf("User: %v\n", user)
	// t.Run("CreateUser", func(t *testing.T) {
	// 	t.Run("Success", func(t *testing.T) {
	// 		// TODO
	// 	})
	// 	t.Run("Failure", func(t *testing.T) {
	// 		// TODO
	// 	})
	// })
	// t.Run("GetUserByEmail", func(t *testing.T) {
	// 	t.Run("Success", func(t *testing.T) {
	// 		// TODO
	// 	})
	// 	t.Run("Failure", func(t *testing.T) {
	// 		// TODO
	// 	})
	// })
}
