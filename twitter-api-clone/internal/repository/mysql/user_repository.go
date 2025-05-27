package mysqlrepo

import (
	"database/sql"
	"errors"
	"twitter-api-clone/core/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.Id, user.Name)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(userId string) error {
	result, err := r.DB.Exec("DELETE FROM users WHERE id = ?", userId)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		// TODO: Melhorar forma de retornar o erro
		return errors.New("user not found")
	}

	return nil
}
