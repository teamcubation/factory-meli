package repositories

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/infrastructure/db/mongo"
	"context"
)

type TweetRepositoryImpl struct{}

func NewTweetRepo() *TweetRepositoryImpl {
	return &TweetRepositoryImpl{}
}

func (t *TweetRepositoryImpl) SaveTweet(user *models.Tweet) error {
	collection, err := mongo.ConnectToCollection("tweets")
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (t *TweetRepositoryImpl) FindAllByUser(id int) ([]*models.Tweet, error) {
	collection, err := mongo.ConnectToCollection("tweets")
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.TODO(), map[string]interface{}{"author.ID": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tweets []*models.Tweet
	for cursor.Next(context.TODO()) {
		var tweet models.Tweet
		if err := cursor.Decode(&tweet); err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}

	return tweets, nil
}
