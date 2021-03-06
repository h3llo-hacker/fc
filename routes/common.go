package routes

import (
	"github.com/h3llo-hacker/fc/handler/challenge"
	"github.com/h3llo-hacker/fc/handler/template"
	"github.com/h3llo-hacker/fc/handler/user"
	"github.com/h3llo-hacker/fc/utils/db"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
	})
}

func validateAuth(auth string) bool {
	C := "auth"
	selector := bson.M{}
	auths, err := db.MongoFind(C, nil, selector)
	if err != nil {
		return false
	}
	for _, Auth := range auths {
		if auth == Auth.(bson.M)["token"] {
			return true
		}
	}
	return false
}

func validateUser(uid string) bool {
	u := user.User{
		UserID: uid,
	}
	return u.UserExist()
}

func validateTemplate(templateID string) bool {
	return template.TemplateExist(templateID)
}

func validateChallenge(challengeID string) bool {
	return challenge.ChallengeExist(challengeID)
}
