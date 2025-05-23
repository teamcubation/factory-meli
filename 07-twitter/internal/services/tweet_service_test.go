package services

import (
	"errors"
	"testing"

	"07-twitter/core/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTweet_Success(t *testing.T) {
	tweetRepo := &mockTweetRepository{
		saveFunc: func(t *models.Tweet) (*models.Tweet, error) {
			return t, nil
		},
	}

	userRepo := &mockUserRepo{}
	service := NewTweetService(tweetRepo, userRepo)
	userID := uuid.New()
	content := "Hello, Twitter!"

	tweet, err := service.CreateTweet(content, userID)

	assert.NoError(t, err)
	assert.Equal(t, content, tweet.Content)
	assert.Equal(t, userID, tweet.UserID)
	assert.NotEqual(t, uuid.Nil, tweet.ID)
}

func Test_CreateTweet_EmptyContent(t *testing.T) {
	service := NewTweetService(&mockTweetRepository{}, &mockUserRepo{})

	_, err := service.CreateTweet("", uuid.New())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conteúdo do tweet não pode ser vazio")
}

func Test_CreateTweet_EmptyUserID(t *testing.T) {
	service := NewTweetService(&mockTweetRepository{}, &mockUserRepo{})

	_, err := service.CreateTweet("algum conteúdo", uuid.Nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID é obrigatório")
}

func Test_CreateTweet_SaveError(t *testing.T) {
	tweetRepo := &mockTweetRepository{
		saveFunc: func(t *models.Tweet) (*models.Tweet, error) {
			return nil, errors.New("erro ao salvar tweet")
		},
	}

	userRepo := &mockUserRepo{}
	service := NewTweetService(tweetRepo, userRepo)
	userID := uuid.New()

	_, err := service.CreateTweet("conteúdo válido", userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao salvar tweet")
}

func Test_GetTweetByUserId_Success(t *testing.T) {
	userID := uuid.New()
	tweets := []*models.Tweet{{ID: uuid.New(), UserID: userID, Content: "tweet1"}}
	tweetRepo := &mockTweetRepository{
		findAllByUserIDFunc: func(id uuid.UUID) ([]*models.Tweet, error) {
			return tweets, nil
		},
	}

	service := NewTweetService(tweetRepo, &mockUserRepo{})
	result, err := service.GetTweetByUserId(userID)

	assert.NoError(t, err)
	assert.Equal(t, tweets, result)
}

func Test_GetTweetByUserId_EmptyUserID(t *testing.T) {
	service := NewTweetService(&mockTweetRepository{}, &mockUserRepo{})

	_, err := service.GetTweetByUserId(uuid.Nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID é obrigatório")
}

func Test_GetTimelineByUserId_Success(t *testing.T) {
	userID := uuid.New()
	followingIDs := []uuid.UUID{uuid.New(), uuid.New()}
	tweets := []*models.Tweet{{ID: uuid.New(), UserID: followingIDs[0], Content: "tweet1"}}
	userRepo := &mockUserRepo{
		findFollowingIDsByUserIDFunc: func(id uuid.UUID) ([]uuid.UUID, error) {
			return followingIDs, nil
		},
	}
	tweetRepo := &mockTweetRepository{
		findTimelineByUsersIdsFunc: func(ids []uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
			return tweets, nil
		},
	}

	service := NewTweetService(tweetRepo, userRepo)
	result, err := service.GetTimelineByUserId(userID, 5, 0)

	assert.NoError(t, err)
	assert.Equal(t, tweets, result)
}

func Test_GetTimelineByUserId_EmptyUserID(t *testing.T) {
	service := NewTweetService(&mockTweetRepository{}, &mockUserRepo{})

	_, err := service.GetTimelineByUserId(uuid.Nil, 5, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID é obrigatório")
}

func Test_GetTimelineByUserId_NoFollowing(t *testing.T) {
	userID := uuid.New()
	userRepo := &mockUserRepo{
		findFollowingIDsByUserIDFunc: func(id uuid.UUID) ([]uuid.UUID, error) {
			return []uuid.UUID{}, nil
		},
	}

	service := NewTweetService(&mockTweetRepository{}, userRepo)
	result, err := service.GetTimelineByUserId(userID, 5, 0)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func Test_GetTimelineByUserId_UserRepoError(t *testing.T) {
	userID := uuid.New()
	userRepo := &mockUserRepo{
		findFollowingIDsByUserIDFunc: func(id uuid.UUID) ([]uuid.UUID, error) {
			return nil, errors.New("repo error")
		},
	}

	service := NewTweetService(&mockTweetRepository{}, userRepo)
	_, err := service.GetTimelineByUserId(userID, 5, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo error")
}

func Test_GetTimelineByUserId_TweetRepoError(t *testing.T) {
	userID := uuid.New()
	followingIDs := []uuid.UUID{uuid.New()}
	userRepo := &mockUserRepo{
		findFollowingIDsByUserIDFunc: func(id uuid.UUID) ([]uuid.UUID, error) {
			return followingIDs, nil
		},
	}
	tweetRepo := &mockTweetRepository{
		findTimelineByUsersIdsFunc: func(ids []uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
			return nil, errors.New("repo error")
		},
	}

	service := NewTweetService(tweetRepo, userRepo)
	_, err := service.GetTimelineByUserId(userID, 5, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo error")
}
