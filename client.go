package main

import (
	"context"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/user/follow"
	followingtypes "github.com/michimani/gotwi/user/follow/types"
	"github.com/michimani/gotwi/user/userlookup"
	userlookuptypes "github.com/michimani/gotwi/user/userlookup/types"
)

func GetUserID(username string, c *gotwi.Client) string {
	p := &userlookuptypes.GetByUsernameInput{
		Username: username,
		Expansions: fields.ExpansionList{
			fields.ExpansionPinnedTweetID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
		},
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldCreatedAt,
		},
	}
	user, err := userlookup.GetByUsername(context.Background(), c, p)
	if err != nil {
		panic(err)
	}
	return *user.Data.ID
}

func GetFollowings(userID, accessToken string, c *gotwi.Client) []resources.User {
	followings := []resources.User{}
	paginationToken := ""
	for {
		p := &followingtypes.ListFollowingsInput{
			ID:              userID,
			MaxResults:      1000,
			PaginationToken: paginationToken,
		}
		p.SetAccessToken(accessToken)

		//unfollow user
		l, err := follow.ListFollowings(context.Background(), c, p)
		if err != nil {
			panic(err)
		}
		followings = append(followings, l.Data...)

		if l.Meta.NextToken == nil {
			break
		}
		paginationToken = *l.Meta.NextToken
	}

	return followings
}
