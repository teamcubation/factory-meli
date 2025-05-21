package out

import "Yuri/twitter/core/models"

type IUserRepository interface {
	SaveUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetById(id string) (*models.User, error)
}
