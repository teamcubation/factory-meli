package mappers

import (
	"Yuri/twitter/core/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tweetMapper struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	AuthorId    string             `bson:"author_id"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func TweetToMongo(tweet *models.Tweet) (*tweetMapper, error) {
	var oid primitive.ObjectID
	var err error
	if tweet.ID == "" {
		oid = primitive.NewObjectID()
	} else {
		oid, err = primitive.ObjectIDFromHex(tweet.ID)
		if err != nil {
			return nil, err
		}
	}
	return &tweetMapper{
		ID:          oid,
		Title:       tweet.Title,
		Description: tweet.Description,
		AuthorId:    tweet.AuthorId,
		CreatedAt:   tweet.CreatedAt,
	}, nil
}

func MongoToTweet(mapper *tweetMapper) *models.Tweet {
	return &models.Tweet{
		ID:          mapper.ID.Hex(),
		Title:       mapper.Title,
		Description: mapper.Description,
		AuthorId:    mapper.AuthorId,
		CreatedAt:   mapper.CreatedAt,
	}
}

func ParamToFilter(key, value string) *bson.M {
	return &bson.M{
		key: value,
	}
}

func ParamToArrayFilter(key string, value []string) *bson.M {
	return &bson.M{
		key: bson.M{"$in": value},
	}
}

func NewOptionsPagination(page, tweetsPerPage int) *options.FindOptions {
	skip := int64((page - 1) * tweetsPerPage)
	options := options.Find()
	options.SetSort(bson.D{{Key: "created_at", Value: -1}})
	options.SetLimit(int64(tweetsPerPage))
	options.SetSkip(skip)
	return options
}
