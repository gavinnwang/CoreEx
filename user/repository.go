package user

import (
	"database/sql"
	"github/wry-0313/exchange/models"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(user models.User) error

	GetUser(userID uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	DeleteUser(userID uuid.UUID) error
}

type repository struct {
	db   *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user models.User) error {

	_, err := r.db.Exec("INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}