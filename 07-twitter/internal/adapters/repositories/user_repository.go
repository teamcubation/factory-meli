package repositories

import (
	model "07-twitter/core/models"
	repo "07-twitter/core/ports/repos"
	"database/sql"
	"errors"

	"github.com/google/uuid"
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

// FindById() retorna um usuário pelo seu id.
func (r *userRepositoryImpl) FindById(id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRow(`SELECT * FROM users WHERE id=$1`, id.String())

	var user model.User
	var idStr string

	if err := row.Scan(&idStr, &user.Username); err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("[uuid.Parse()] - " + err.Error())
	}

	rows, err := r.db.Query(`SELECT follow_id FROM follows WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}

	following := []uuid.UUID{}
	for rows.Next() {
		if err := rows.Scan(&following); err != nil {
			return nil, err
		}
	}

	user.Following = following
	user.ID = parsedID
	return &user, nil
}
