package mysqlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"twitter-api-clone/core/models"
)

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) Follow(follow models.Follow) error {
	fmt.Printf("INSERT INTO follows (follower_id, followed_id, created_at) VALUES (%v,%v,%v)", follow.FollowerId, follow.FollowedId, follow.CreatedAt)

	_, err := r.DB.Exec(
		"INSERT INTO follows (follower_id, followed_id, created_at) VALUES (?,?,?)",
		follow.FollowerId, follow.FollowedId, follow.CreatedAt,
	)

	return err
}

func (r *FollowRepository) Unfollow(follow models.Follow) error {
	result, err := r.DB.Exec("DELETE FROM follows WHERE follower_id = ? AND followed_id = ?", follow.FollowerId, follow.FollowedId)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		// TODO: Melhorar forma de retornar o erro
		return errors.New("follow not found")
	}

	return err
}

// func (r *FollowRepository) ListByUser(userId string) ([]models.Tweet, error) {
// 	rows, err := r.DB.Query("SELECT follower_id, followed_id, created_at FROM follows WHERE user_id = ? ORDER BY created_at DESC", userId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var tweets []models.Tweet
// 	for rows.Next() {
// 		var t models.Tweet
// 		err := rows.Scan(&t.Id, &t.UserId, &t.Content, &t.CreatedAt)

// 		if err != nil {
// 			fmt.Printf("to aqui %v", err.Error())
// 			return nil, err
// 		}

// 		tweets = append(tweets, t)
// 	}

// 	return tweets, nil
// }
