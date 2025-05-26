package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

// TweetInteractor defines the interface for tweet-related actions like liking, unliking, and retweeting.
type TweetInteractorServices interface {
	Like(ctx context.Context, r *http.Request) error
	Unlike(ctx context.Context, r *http.Request) error
	Retweet(ctx context.Context, r *http.Request) error
}

// TweetInteractorServicesImpl implements the TweetInteractor interface.
type TweetInteractorServicesImpl struct {
	u_repository postgres.UserRepository
	l_repository postgres.LikeRepository
	t_repository postgres.TweetRepository
	r_repository postgres.RetweetRepository
}

// NewTweetInteractor creates a new TweetInteractorServicesImpl instance.
func NewTweetInteractor(userRepo postgres.UserRepository, likeRepo postgres.LikeRepository, tweetRepo postgres.TweetRepository, retweetRepo postgres.RetweetRepository) TweetInteractorServices {
	return &TweetInteractorServicesImpl{
		u_repository: userRepo,
		l_repository: likeRepo,
		t_repository: tweetRepo,
		r_repository: retweetRepo,
	}
}

// Like allows a user to like a tweet.
func (s TweetInteractorServicesImpl) Like(ctx context.Context, r *http.Request) error {
	slog.InfoContext(ctx, "Calling service to like a tweet", "layer", "service")

	// Parse request body for UserID
	type params struct {
		UserIDStr string `json:"user_id"`
	}

	requestParams := params{}
	if err := json.NewDecoder(r.Body).Decode(&requestParams); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for Like", "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	userID, err := uuid.Parse(requestParams.UserIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid user ID from payload for Like", "user_id_str", requestParams.UserIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid user ID"}
	}

	tweetIDStr := r.PathValue("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid tweet ID from path for Like", "tweet_id_str", tweetIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid tweet ID"}
	}

	// Verify user exists
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{ID: userID})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for Like", "user_id", userID, "layer", "service")
			return ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for Like", "user_id", userID, "error", err, "layer", "service")
		return err
	}

	// Verify tweet exists
	_, err = s.t_repository.FindByID(ctx, postgres.TweetFindByIDParams{ID: tweetID}) // Assuming TweetRepository has a FindByID
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Tweet not found for Like", "tweet_id", tweetID, "layer", "service")
			return ServiceError{http.StatusNotFound, "tweet not found"}
		}
		slog.ErrorContext(ctx, "Error finding tweet by ID for Like", "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}

	// Check if already liked
	hasLiked, err := s.l_repository.HasLike(ctx, postgres.LikeHasParams{UserID: userID, TweetID: tweetID})
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if user has liked tweet", "user_id", userID, "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}
	if hasLiked {
		slog.WarnContext(ctx, "Tweet already liked by user", "user_id", userID, "tweet_id", tweetID, "layer", "service")
		return ServiceError{http.StatusConflict, "tweet already liked by this user"}
	}

	// Add the like
	_, err = s.l_repository.AddLike(ctx, postgres.LikeAddParams{UserID: userID, TweetID: tweetID})
	if err != nil {
		slog.ErrorContext(ctx, "Error adding like", "user_id", userID, "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}

	slog.InfoContext(ctx, "Successfully liked tweet", "user_id", userID, "tweet_id", tweetID, "layer", "service")
	return nil
}

func (s TweetInteractorServicesImpl) Unlike(ctx context.Context, r *http.Request) error {
	slog.InfoContext(ctx, "Calling service to unlike a tweet", "layer", "service")

	// Parse request body for UserID
	type params struct {
		UserIDStr string `json:"user_id"`
	}

	requestParams := params{}
	if err := json.NewDecoder(r.Body).Decode(&requestParams); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for Unlike", "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	userID, err := uuid.Parse(requestParams.UserIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid user ID from payload for Unlike", "user_id_str", requestParams.UserIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid user ID"}
	}

	tweetIDStr := r.PathValue("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid tweet ID from path for Unlike", "tweet_id_str", tweetIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid tweet ID"}
	}

	// Verify user exists
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{ID: userID})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for Unlike", "user_id", userID, "layer", "service")
			return ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for Unlike", "user_id", userID, "error", err, "layer", "service")
		return err
	}

	// Verify tweet exists
	_, err = s.t_repository.FindByID(ctx, postgres.TweetFindByIDParams{ID: tweetID}) 
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Tweet not found for Unlike", "tweet_id", tweetID, "layer", "service")
			return ServiceError{http.StatusNotFound, "tweet not found"}
		}
		slog.ErrorContext(ctx, "Error finding tweet by ID for Unlike", "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}

	// Check if the like exists before attempting to remove
	hasLiked, err := s.l_repository.HasLike(ctx, postgres.LikeHasParams{UserID: userID, TweetID: tweetID})
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if user has liked tweet before unliking", "user_id", userID, "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}
	if !hasLiked {
		slog.WarnContext(ctx, "Tweet not liked by user, cannot unlike", "user_id", userID, "tweet_id", tweetID, "layer", "service")
		return ServiceError{http.StatusNotFound, "tweet not liked by this user"}
	}

	// Remove the like
	err = s.l_repository.RemoveLike(ctx, postgres.LikeRemoveParams{UserID: userID, TweetID: tweetID})
	if err != nil {
		slog.ErrorContext(ctx, "Error removing like", "user_id", userID, "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}

	slog.InfoContext(ctx, "Successfully unliked tweet", "user_id", userID, "tweet_id", tweetID, "layer", "service")
	return nil
}

// Retweet allows a user to retweet a tweet.
func (s TweetInteractorServicesImpl) Retweet(ctx context.Context, r *http.Request) error {
	slog.InfoContext(ctx, "Calling service to retweet a tweet", "layer", "service")

	// Parse request body for UserID
	type params struct {
		UserIDStr string `json:"user_id"`
	}

	requestParams := params{}
	if err := json.NewDecoder(r.Body).Decode(&requestParams); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for Retweet", "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	userID, err := uuid.Parse(requestParams.UserIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid user ID from payload for Retweet", "user_id_str", requestParams.UserIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid user ID"}
	}

	tweetIDStr := r.PathValue("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid tweet ID from path for Retweet", "tweet_id_str", tweetIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid tweet ID"}
	}

	// Verify user exists
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{ID: userID})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for Retweet", "user_id", userID, "layer", "service")
			return ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for Retweet", "user_id", userID, "error", err, "layer", "service")
		return err
	}

	// Verify original tweet exists and get its creator ID
	originalTweet, err := s.t_repository.FindByID(ctx, postgres.TweetFindByIDParams{ID: tweetID}) // Assuming TweetRepository has a FindByID
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Original tweet not found for Retweet", "tweet_id", tweetID, "layer", "service")
			return ServiceError{http.StatusNotFound, "original tweet not found"}
		}
		slog.ErrorContext(ctx, "Error finding original tweet by ID for Retweet", "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}

	// Prevent a user from retweeting their own tweet
	if originalTweet.CreatorID == userID {
		slog.WarnContext(ctx, "User cannot retweet their own tweet", "user_id", userID, "tweet_id", tweetID, "layer", "service")
		return ServiceError{http.StatusBadRequest, "cannot retweet your own tweet"}
	}

	// Check if already retweeted
	hasRetweeted, err := s.r_repository.HasRetweet(ctx, postgres.RetweetHasParams{UserID: userID, TweetID: tweetID})
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if user has retweeted tweet", "user_id", userID, "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}
	if hasRetweeted {
		slog.WarnContext(ctx, "Tweet already retweeted by user", "user_id", userID, "tweet_id", tweetID, "layer", "service")
		return ServiceError{http.StatusConflict, "tweet already retweeted by this user"}
	}

	// Add the retweet
	_, err = s.r_repository.AddRetweet(ctx, postgres.RetweetAddParams{
		UserID:          userID,
		TweetID:         tweetID,
		OriginalTweetID: tweetID, // A retweet directly references the original tweet
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error adding retweet", "user_id", userID, "tweet_id", tweetID, "error", err, "layer", "service")
		return err
	}

	slog.InfoContext(ctx, "Successfully retweeted tweet", "user_id", userID, "tweet_id", tweetID, "layer", "service")
	return nil
}
