package services

import (
	"database/sql"
	"errors"
	"testing"

	model "07-twitter/core/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	repo := &mockUserRepo{
		saveFunc: func(user *model.User) error { return nil },
	}
	service := NewUserService(repo)

	user, err := service.CreateUser("testuser")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestCreateUser_Invalid(t *testing.T) {
	repo := &mockUserRepo{}
	service := NewUserService(repo)

	_, err := service.CreateUser("")

	assert.Error(t, err)
}

func TestCreateUser_SaveError(t *testing.T) {
	repo := &mockUserRepo{
		saveFunc: func(user *model.User) error { return errors.New("erro ao salvar") },
	}
	service := NewUserService(repo)

	_, err := service.CreateUser("testuser")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo.Save")
}

func TestGetUserById_Success(t *testing.T) {
	id := uuid.New()
	repo := &mockUserRepo{
		findByIdFunc: func(uid uuid.UUID) (*model.User, error) {
			return &model.User{ID: uid, Username: "testuser"}, nil
		},
	}
	service := NewUserService(repo)

	user, err := service.GetUserById(id.String())

	assert.NoError(t, err)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetUserById_InvalidUUID(t *testing.T) {
	repo := &mockUserRepo{}
	service := NewUserService(repo)

	_, err := service.GetUserById("not-a-uuid")

	assert.Error(t, err)
}

func TestGetUserById_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &mockUserRepo{
		findByIdFunc: func(uid uuid.UUID) (*model.User, error) {
			return nil, sql.ErrNoRows
		},
	}
	service := NewUserService(repo)

	_, err := service.GetUserById(id.String())

	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestGetUserById_RepoError(t *testing.T) {
	id := uuid.New()
	repo := &mockUserRepo{
		findByIdFunc: func(uid uuid.UUID) (*model.User, error) {
			return nil, errors.New("erro repo")
		},
	}
	service := NewUserService(repo)

	_, err := service.GetUserById(id.String())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo.FindById")
}

func TestFollow_Success(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		followFunc: func(userId, followingId uuid.UUID) error {
			called = true
			return nil
		},
	}
	service := NewUserService(repo)

	err := service.Follow(uuid.New(), uuid.New())

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestFollow_Error(t *testing.T) {
	repo := &mockUserRepo{
		followFunc: func(userId, followingId uuid.UUID) error {
			return errors.New("erro follow")
		},
	}
	service := NewUserService(repo)

	err := service.Follow(uuid.New(), uuid.New())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro follow")
}

func TestUnfollow_Success(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		unfollowFunc: func(userId, followingId uuid.UUID) error {
			called = true
			return nil
		},
	}
	service := NewUserService(repo)

	err := service.Unfollow(uuid.New(), uuid.New())

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestUnfollow_Error(t *testing.T) {
	repo := &mockUserRepo{
		unfollowFunc: func(userId, followingId uuid.UUID) error {
			return errors.New("erro unfollow")
		},
	}
	service := NewUserService(repo)

	err := service.Unfollow(uuid.New(), uuid.New())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro unfollow")
}
