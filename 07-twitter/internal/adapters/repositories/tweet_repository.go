package repositories

import (
	"07-twitter/core/models"
	repo "07-twitter/core/ports/repos"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// tweetRepositoryImpl implementa a interface TweetRepository usando um banco SQL.
type tweetRepositoryImpl struct {
	db *sql.DB
}

// NewTweetRepository() cria uma nova instância de TweetRepository com a conexão fornecida.
func NewTweetRepository(db *sql.DB) repo.TweetRepository {
	return &tweetRepositoryImpl{db: db}
}

// Save() insere um novo tweet no banco de dados e retorna o tweet salvo com os campos preenchidos.
func (r *tweetRepositoryImpl) Save(tweet *models.Tweet) (*models.Tweet, error) {
	var savedTweet models.Tweet

	err := r.db.QueryRow(
		`INSERT INTO tweets (id, user_id, content) VALUES ($1, $2, $3) RETURNING id, user_id, content, created_at`,
		tweet.ID,
		tweet.UserID,
		tweet.Content,
	).Scan(&savedTweet.ID, &savedTweet.UserID, &savedTweet.Content, &savedTweet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &savedTweet, nil
}

// FindAllByUserID() retorna todos os tweets de um usuário específico, ordenados do mais recente para o mais antigo.
func (r *tweetRepositoryImpl) FindAllByUserID(userID uuid.UUID) ([]*models.Tweet, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, content, created_at FROM tweets WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	var tweets []*models.Tweet
	for rows.Next() {
		var tweet models.Tweet

		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt); err != nil {
			return nil, err
		}

		tweets = append(tweets, &tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tweets, nil
}

// FindTimelineByUsersIds() retorna os tweets dos usuários informados, ordenados do mais recente para o mais antigo, com paginação.
func (r *tweetRepositoryImpl) FindTimelineByUsersIds(userIDs []uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
	if len(userIDs) == 0 {
		return []*models.Tweet{}, nil
	}

	query := `SELECT id, user_id, content, created_at FROM tweets WHERE user_id IN (`
	args := make([]interface{}, len(userIDs)+2)
	for i, id := range userIDs {
		query += "$" + fmt.Sprintf("%d", i+1)
		if i < len(userIDs)-1 {
			query += ","
		}
		args[i] = id
	}
	query += `) ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(userIDs)+1) +
		` OFFSET $` + fmt.Sprintf("%d", len(userIDs)+2)

	args[len(userIDs)] = limit
	args[len(userIDs)+1] = offset

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var tweets []*models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tweets, nil
}
