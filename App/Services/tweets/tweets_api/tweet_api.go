package tweets_api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ahmed-AbdElRhman/twitter-architecture/tweets/tweets_services"
	"github.com/labstack/echo/v4"
)

type TweetManager interface {
	CreateTweet(tweet *tweets_services.Tweet) error
	GetUserTweets(ownerId *int) ([]*tweets_services.Tweet, error)
}

type TweetRouter struct {
	tweetMngr TweetManager
}

func NewTweetRouter(tweetMngr TweetManager) *TweetRouter {
	return &TweetRouter{
		tweetMngr: tweetMngr,
	}
}

func (obj *TweetRouter) CreateTweet(c echo.Context) error {
	// Get the request body
	var tweet tweets_services.Tweet
	err := c.Bind(&tweet)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	//------ Create Tweet Services Logic -------------
	err = obj.tweetMngr.CreateTweet(&tweet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Tweet Created")
}
func (obj *TweetRouter) GetUserTweets(c echo.Context) error {
	// Get the request body
	ownerId, err := strconv.Atoi(c.QueryParam("ownerId"))
	fmt.Println("ownerId", ownerId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid ownerId"))
	}
	if ownerId == 0 {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid ownerId"))
	}
	//------ Get User Tweets Services Logic -------------
	tweets, err := obj.tweetMngr.GetUserTweets(&ownerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tweets)
}
