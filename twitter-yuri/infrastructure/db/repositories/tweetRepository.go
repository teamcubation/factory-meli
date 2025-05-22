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

func (t *TweetRepositoryImpl) FindAllByUser(userID string, page, tweetsPerPage int) ([]*models.Tweet, error) {
	collection, err := mongo.ConnectToCollection("tweets")
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(context.TODO(), mappers.ParamToFilter("author_id", userID), mappers.NewOptionsPagination(page, tweetsPerPage))
	if err != nil {
		return nil, err
	}

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
	defer cursor.Close(context.TODO())
	if errDisconnect := mongo.Disconnect(collection.Database().Client()); errDisconnect != nil {
		return nil, errDisconnect
	}
	return tweets, nil
}

func (t *TweetRepositoryImpl) FindAllByUserTimeline(userIds []string, page int, tweetsPerPage int) ([]*models.Tweet, error) {
	collection, err := mongo.ConnectToCollection("tweets")
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(context.TODO(), mappers.ParamToArrayFilter("author_id", userIds), mappers.NewOptionsPagination(page, tweetsPerPage))
	if err != nil {
		return nil, err
	}

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
	defer cursor.Close(context.TODO())
	if errDisconnect := mongo.Disconnect(collection.Database().Client()); errDisconnect != nil {
		return nil, errDisconnect
	}
	return tweets, nil
}
