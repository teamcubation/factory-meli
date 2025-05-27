package mysqlrepo

import (
	"database/sql"
	"fmt"
	"twitter-api-clone/core/models"
)

type TweetRepository struct {
	DB *sql.DB
}

func NewTweetRepository(db *sql.DB) *TweetRepository {
	return &TweetRepository{DB: db}
}

func (r *TweetRepository) Create(tweet *models.Tweet) error {
	_, err := r.DB.Exec(
		"INSERT INTO tweets (id, user_id, content, created_at) VALUES (?,?,?,?)",
		tweet.Id, tweet.UserId, tweet.Content, tweet.CreatedAt,
	)

	return err
}

func (r *TweetRepository) ListByUser(userId string) ([]models.Tweet, error) {
	rows, err := r.DB.Query("SELECT id, user_id, content, created_at FROM tweets WHERE user_id = ? ORDER BY created_at DESC", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var t models.Tweet
		err := rows.Scan(&t.Id, &t.UserId, &t.Content, &t.CreatedAt)

		if err != nil {
			fmt.Printf("to aqui %v", err.Error())
			return nil, err
		}

		tweets = append(tweets, t)
	}

	return tweets, nil
}

func (r *TweetRepository) GetTimeline(userId string, limit int, offset int) ([]models.Timeline, error) {
	rows, err := r.DB.Query(`
		SELECT 
			t.id as tweet_id,
			u.id as user_id,
			u.name as user_name,
			t.content  as tweet_content,
			t.created_at as tweet_date
		FROM tweets t
		INNER JOIN follows f
		on t.user_id = f.followed_id
		INNER JOIN users u 
		ON u.id = t.user_id 
		WHERE 1=1
		AND f.follower_id = ?
		ORDER BY tweet_date DESC
		LIMIT ? OFFSET ?
	`, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var timeline []models.Timeline

	for rows.Next() {
		var t models.Timeline
		err := rows.Scan(&t.TweetId, &t.UserId, &t.Name, &t.Content, &t.CreatedAt)

		if err != nil {
			fmt.Printf("to aqui %v", err.Error())
			return nil, err
		}

		timeline = append(timeline, t)
	}

	return timeline, nil
}
