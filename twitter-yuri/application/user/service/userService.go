package user

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
	"fmt"
	"slices"
)

type UserServiceImpl struct {
	repo out.IUserRepository
}

func NewUserService(repo out.IUserRepository) services.IUserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) CreateUser(user *models.User) error {

	if err := s.repo.SaveUser(user); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) FollowUser(idUser, id string) error {
	user, errFollow := s.checkCanFollow(idUser, id)
	if errFollow != nil {
		return errFollow
	}
	var model = models.User{
		ID:          idUser,
		Name:        user.Name,
		FollowUsers: user.FollowUsers,
	}
	model.FollowUsers = append(model.FollowUsers, id)
	if err := s.repo.UpdateUser(&model); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) UnfollowUser(idUser, id string) error {
	user, errUnfollow := s.checkCanUnfollow(idUser, id)
	if errUnfollow != nil {
		return errUnfollow
	}
	var model = models.User{
		ID:          idUser,
		Name:        user.Name,
		FollowUsers: user.FollowUsers,
	}
	model.FollowUsers = slices.DeleteFunc(user.FollowUsers, func(followedID string) bool { return followedID == id })
	// model.FollowUsers = delete(model.FollowUsers,id)
	if err := s.repo.UpdateUser(&model); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) checkCanFollow(idUser, followId string) (*models.User, error) {
	user, errGet := s.GetUserById(idUser)
	if errGet != nil {
		return nil, errGet
	}
	if slices.Contains(user.FollowUsers, followId) {
		return nil, fmt.Errorf("user %s already follows userId %s", user.Name, followId)
	}
	return user, nil
}

func (s *UserServiceImpl) checkCanUnfollow(idUser, followId string) (*models.User, error) {
	user, errGet := s.GetUserById(idUser)
	if errGet != nil {
		return nil, errGet
	}
	if !slices.Contains(user.FollowUsers, followId) {
		return nil, fmt.Errorf("user %s not follow the userId %s", user.Name, followId)
	}
	return user, nil
}
func (s *UserServiceImpl) GetUserById(id string) (*models.User, error) {
	user, errGet := s.repo.GetById(id)
	if user == nil {
		return nil, fmt.Errorf("user %s not found", id)
	}
	if errGet != nil {
		return nil, errGet
	}
	return user, nil
}
