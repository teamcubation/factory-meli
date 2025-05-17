package repositories

import (
	model "07-twitter/core/models"
	repo "07-twitter/core/ports/repos"
	"database/sql"

	// "github.com/google/uuid"
)

// userRepositoryImpl é uma implementação PostgreSQL de UserRepository.
type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository cria uma nova instância de userRepositoryImpl com banco PostgreSQL.
func NewUserRepository(db *sql.DB) repo.UserRepository {
	return &userRepositoryImpl{db: db}
}

// Save armazena um usuário no banco de dados.
func (r *userRepositoryImpl) Save(user *model.User) error {
	_, err := r.db.Exec(`INSERT INTO users (id, username) VALUES ($1, $2)`, user.ID.String(), user.Username)
	return err
}