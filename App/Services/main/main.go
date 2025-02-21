package main

import (
	"fmt"

	"github.com/Ahmed-AbdElRhman/twitter-architecture/hashingapp"
	nosqlmngr "github.com/Ahmed-AbdElRhman/twitter-architecture/nosqlMngr"
	sqlmngr "github.com/Ahmed-AbdElRhman/twitter-architecture/sqlMngr"
	"github.com/Ahmed-AbdElRhman/twitter-architecture/tweets/tweets_api"
	"github.com/Ahmed-AbdElRhman/twitter-architecture/tweets/tweets_services"
	"github.com/Ahmed-AbdElRhman/twitter-architecture/users/authmiddleware"
	"github.com/Ahmed-AbdElRhman/twitter-architecture/users/users_api"
	"github.com/Ahmed-AbdElRhman/twitter-architecture/users/users_services"
	"github.com/Ahmed-AbdElRhman/twitter-architecture/utils"
	"github.com/labstack/echo/v4"
)

// type SqlMngr2 interface {
// 	CreateTables() error
// 	SeedProducts() error
// }

var (
	host     = utils.Host
	port     = utils.Port
	user     = utils.User
	password = utils.Password
	dbname   = utils.DBName
)

func main() {
	// ********************** USER SERVICE **********************
	// define the type of dbConnector
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(connStr)
	sqldbMngr, err := sqlmngr.NewPostgres(connStr, "../sqlMngr/schema.sql")
	if err != nil {
		panic("Error while connect to Users database:" + err.Error())
	}
	defer sqldbMngr.CloseDB()
	// Hashing Manager
	hashMngr := hashingapp.NewHashMngr(utils.HASHING_SECRET)
	/*
		inject the Hashing Manager to the User Service Interface
		inject the Database to the User Service Interface
	*/
	usersMgr := users_services.NewUsersService(sqldbMngr, hashMngr)
	// define the type of JWT
	jwtObj := authmiddleware.NewLocalMiddlewareMngr(utils.JWT_SECRET)
	//inject the UserMnger feature to the User API interface
	userRouter := users_api.NewUsersRouter(usersMgr, jwtObj)

	// ********************** TWEET SERVICE **********************
	tweetDbMngr, err := nosqlmngr.NewMongoDb("mongodb://admin:admin@localhost:27017")
	if err != nil {
		panic("Error while connect to tweetDbMngr database:" + err.Error())
	}
	//inject the Database to the tweet Service Interface
	tweetMgr := tweets_services.NewTweetService(tweetDbMngr)
	//inject the TweetMnger feature to the Tweet API interface
	tweetRouter := tweets_api.NewTweetRouter(tweetMgr)

	// ******* Define echo server Route Table *******
	e := echo.New()
	// ------ User API Routes ---------
	e.POST("/login", userRouter.Login)
	e.POST("/register", userRouter.Register)
	protected := e.Group("")
	protected.Use(jwtObj.JWTMiddleware())
	protected.POST("/gettweets", userRouter.GetUserTweets, jwtObj.GroupAuthorization([]string{"edit_tweets"}))
	// ------ Tweet API Routes ---------
	protected.POST("/createtweet", tweetRouter.CreateTweet)
	protected.GET("/getusertweets", tweetRouter.GetUserTweets)
	e.Logger.Fatal(e.Start(":8080"))
}
