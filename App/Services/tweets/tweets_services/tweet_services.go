package tweets_services

type TweetDBmngr interface {
	CreateTweet(tweet *Tweet, dbname, collectionName string) error
	GetUserTweets(ownerId *int, dbname, collectionName string) ([]*Tweet, error)
	CloseDB() error
}
type TweetService struct {
	dbMngr TweetDBmngr
}

func NewTweetService(dbMngr TweetDBmngr) *TweetService {
	return &TweetService{
		dbMngr: dbMngr,
	}
}

func (obj *TweetService) CreateTweet(tweet *Tweet) error {
	return obj.dbMngr.CreateTweet(tweet, "tweetsdb", "tweets")
}

func (obj *TweetService) GetUserTweets(ownerId *int) ([]*Tweet, error) {
	return obj.dbMngr.GetUserTweets(ownerId, "tweetsdb", "tweets")
}
