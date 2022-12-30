package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweet/timeline"
	timelinetypes "github.com/michimani/gotwi/tweet/timeline/types"
	"github.com/michimani/gotwi/user/follow"
	followtypes "github.com/michimani/gotwi/user/follow/types"
)

func main() {
	godotenv.Load(".env")
	accessToken := os.Getenv("ACCESS_TOKEN")
	fileToLoad := os.Getenv("FILE_TO_LOAD")

	twitterClient, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           os.Getenv("OAUTH_TOKEN"),
		OAuthTokenSecret:     os.Getenv("OAUTH_TOKEN_SECRET"),
	})

	if err != nil {
		panic(err)
	}

	followings := []resources.User{}
	userID := GetUserID(os.Getenv("USERNAME"), twitterClient)
	if fileToLoad == "" {
		followings = GetFollowings(userID, accessToken, twitterClient)
	} else {
		followings = LoadUsers(fileToLoad)
	}

	for i, v := range followings {
		t, err := timeline.ListTweets(context.Background(), twitterClient, &timelinetypes.ListTweetsInput{
			ID:         *v.ID,
			MaxResults: 5,
			TweetFields: fields.TweetFieldList{
				fields.TweetFieldCreatedAt,
			},
		})
		if err != nil {
			panic(err)
		}

		// if CreatedAt more then 3 months ago, unfollow
		if len(t.Data) == 0 || t.Data[0].CreatedAt.Before(time.Now().AddDate(0, -3, 0)) {
			fmt.Printf("%d/%d: Try to Unfollow %s", i+1, len(followings), *v.Username)
			p := &followtypes.DeleteFollowingInput{
				SourceUserID: userID,
				TargetID:     *v.ID,
			}
			p.SetAccessToken(accessToken)
			_, err := follow.DeleteFollowing(context.Background(), twitterClient, p)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("User %s deleted! \n", *v.Username)
			}
		}

		fmt.Printf("%d/%d\r", i, len(followings))
		time.Sleep(20 * time.Second)
	}
}
