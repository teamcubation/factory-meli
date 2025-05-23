package services

import (
	model "07-twitter/core/models"

	"github.com/google/uuid"
)

type mockTweetRepository struct {
	saveFunc                   func(*model.Tweet) (*model.Tweet, error)
	findAllByUserIDFunc        func(uuid.UUID) ([]*model.Tweet, error)
	findTimelineByUsersIdsFunc func([]uuid.UUID, int, int) ([]*model.Tweet, error)
}

func (m *mockTweetRepository) Save(t *model.Tweet) (*model.Tweet, error) {
	return m.saveFunc(t)
}

func (m *mockTweetRepository) FindAllByUserID(id uuid.UUID) ([]*model.Tweet, error) {
	if m.findAllByUserIDFunc != nil {
		return m.findAllByUserIDFunc(id)
	}
	return nil, nil
}

func (m *mockTweetRepository) FindTimelineByUsersIds(ids []uuid.UUID, limit, offset int) ([]*model.Tweet, error) {
	if m.findTimelineByUsersIdsFunc != nil {
		return m.findTimelineByUsersIdsFunc(ids, limit, offset)
	}
	return nil, nil
}

type mockUserRepo struct {
	saveFunc                     func(*model.User) error
	findByIdFunc                 func(uuid.UUID) (*model.User, error)
	followFunc                   func(uuid.UUID, uuid.UUID) error
	unfollowFunc                 func(uuid.UUID, uuid.UUID) error
	findFollowingIDsByUserIDFunc func(uuid.UUID) ([]uuid.UUID, error)
}

func (m *mockUserRepo) Save(user *model.User) error {
	if m.saveFunc != nil {
		return m.saveFunc(user)
	}
	return nil
}

func (m *mockUserRepo) FindById(id uuid.UUID) (*model.User, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	return nil, nil
}

func (m *mockUserRepo) Follow(userId, followingId uuid.UUID) error {
	if m.followFunc != nil {
		return m.followFunc(userId, followingId)
	}
	return nil
}

func (m *mockUserRepo) Unfollow(userId, followingId uuid.UUID) error {
	if m.unfollowFunc != nil {
		return m.unfollowFunc(userId, followingId)
	}
	return nil
}
func (m *mockUserRepo) FindFollowingIDsByUserID(id uuid.UUID) ([]uuid.UUID, error) {
	if m.findFollowingIDsByUserIDFunc != nil {
		return m.findFollowingIDsByUserIDFunc(id)
	}
	return nil, nil
}
