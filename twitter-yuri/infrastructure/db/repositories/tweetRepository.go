package repositories

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/infrastructure/db/mongo"
	"Yuri/twitter/infrastructure/db/mongo/mappers"
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
	mongoTweet, err := mappers.TweetToMongo(user)
	if err != nil {
		return err
	}

	if _, errInsert := collection.InsertOne(context.TODO(), mongoTweet); errInsert != nil {
		return err
	}

	if errDisconnect := mongo.Disconnect(collection.Database().Client()); errDisconnect != nil {
		return errDisconnect
	}
	return nil
}

func (t *TweetRepositoryImpl) FindAllByUser(userID string) ([]*models.Tweet, error) {
	collection, err := mongo.ConnectToCollection("tweets")
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.TODO(), map[string]interface{}{"author_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tweets []*models.Tweet
	for cursor.Next(context.TODO()) {
		var tweet models.Tweet
		tweetMapper, errMap := mappers.TweetToMongo(&tweet)
		if errMap != nil {
			return nil, errMap
		}
		if err := cursor.Decode(tweetMapper); err != nil {
			return nil, err
		}
		tweets = append(tweets, mappers.MongoToTweet(tweetMapper))
	}
	return tweets, nil
}
