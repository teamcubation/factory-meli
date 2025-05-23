package databases

import (
	"errors"
	"fmt"
	"time"
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/models"
	"tweet-tq-rafael/utils"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB() (*gorm.DB, error) {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "admin")
	dbPassword := utils.GetEnv("DB_PASSWORD", "1234")
	dbName := utils.GetEnv("DB_NAME", "twitter-tq")

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Tweet{})
	db.AutoMigrate(&models.Follow{})

	return db, nil
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(createUserDTO dtos.CreateUserRequestDTO) (models.User, error) {

	user := models.User{
		ID:        uuid.New(),
		Name:      createUserDTO.Name,
		Email:     createUserDTO.Email,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	r.db.Create(&user)

	return user, nil
}

func (r *UserRepository) Follow(followRequestDTO dtos.FollowRequestDTO) error {

	var followerUser models.User
	if err := r.db.First(&followerUser, followRequestDTO.FollowerID).Error; err != nil {
		return errors.New("follower user doesn't exist")
	}

	var followedUser models.User
	if err := r.db.First(&followedUser, followRequestDTO.FollowedID).Error; err != nil {
		return errors.New("followed user doesn't exist")
	}

	var existing models.Follow
	if err := r.db.
		Where("follower_id = ? AND followed_id = ?", followRequestDTO.FollowerID, followRequestDTO.FollowedID).
		First(&existing).Error; err == nil {
		return errors.New("already following this user")
	}

	follow := models.Follow{FollowerID: followRequestDTO.FollowerID, FollowedID: followRequestDTO.FollowedID}
	if err := r.db.Create(&follow).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Unfollow(unfollowRequestDTO dtos.UnfollowRequestDTO) error {

	var followerUser models.User
	if err := r.db.First(&followerUser, unfollowRequestDTO.FollowerID).Error; err != nil {
		return errors.New("follower user doesn't exit")
	}

	var followedUser models.User
	if err := r.db.First(&followedUser, unfollowRequestDTO.FollowedID).Error; err != nil {
		return errors.New("followed user doesn't exist")
	}

	var existing models.Follow
	if err := r.db.Where("follower_id = ? AND followed_id = ?", unfollowRequestDTO.FollowerID, unfollowRequestDTO.FollowedID).
		First(&existing).Error; err != nil {
		return nil
	}

	if err := r.db.Delete(&existing).Error; err != nil {
		return errors.New("failed to unfollow user")
	}

	return nil
}

func (r *UserRepository) Timeline(timelineRequestDTO dtos.TimelineRequestDTO) ([]models.Tweet, error) {

	var followedIDs []uuid.UUID

	if err := r.db.
		Model(&models.Follow{}).
		Where("follower_id = ?", timelineRequestDTO.UserID).
		Pluck("followed_id", &followedIDs).Error; err != nil {
		return []models.Tweet{}, err
	}

	var tweets []models.Tweet
	if err := r.db.
		Preload("User").
		Where("user_id IN ?", followedIDs).
		Order("created_at DESC").
		Find(&tweets).Error; err != nil {
		return []models.Tweet{}, err
	}

	return tweets, nil

}

type TweetRepository struct {
	db *gorm.DB
}

func (r *TweetRepository) Create(createTweetRequestDTO dtos.CreateTweetRequestDTO) (models.Tweet, error) {
	tweet := models.Tweet{
		ID:      uuid.New(),
		Content: createTweetRequestDTO.Content,
		UserID:  createTweetRequestDTO.UserID,
	}

	if err := r.db.Create(&tweet).Error; err != nil {
		return models.Tweet{}, err
	}

	return tweet, nil
}

type Repositories struct {
	UserRepository  *UserRepository
	TweetRepository *TweetRepository
}

func RepositoriesFactory(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository:  &UserRepository{db},
		TweetRepository: &TweetRepository{db},
	}
}
