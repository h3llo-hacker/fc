package routes

import (
	"fmt"
	"handler/challenge"

	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
)

func challenges(c *gin.Context) {
	challenges, err := challenge.AllChallenges()
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, challenges)
	}
}

func challengeInfo(c *gin.Context) {
	cid := c.Param("challengeID")
	challenge, err := challenge.QueryChallenge(cid)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, challenge)
	}
}

func challengeCreate(c *gin.Context) {
	userID := c.Request.PostFormValue("uid")
	templateID := c.Request.PostFormValue("templateID")

	if !validateUser(userID) {
		c.JSON(500, gin.H{
			"err": "user [" + userID + "] not found",
		})
		return
	}

	if !validateTemplate(templateID) {
		c.JSON(500, gin.H{
			"err": "template [" + templateID + "] not found",
		})
		return
	}

	uid, _ := uuid.NewV4()
	challengeID := fmt.Sprintf("%v", uid)

	go func() {
		_, err := challenge.CreateChallenge(userID,
			templateID, challengeID)
		if err != nil {
			challenge.UpdateChallengeState(challengeID, "failed")
		} else {
			challenge.UpdateChallengeState(challengeID, "created")
		}
	}()
	c.JSON(201, gin.H{
		"challenge created": "ok",
		"id":                challengeID,
	})
}

func challengeRemove(c *gin.Context) {
	uid := c.Request.PostFormValue("uid")
	challengeID := c.Request.PostFormValue("cid")

	if !validateUser(uid) {
		c.JSON(500, gin.H{
			"err": "uid: [" + uid + "] not found",
		})
		return
	}
	if !validateChallenge(challengeID) {
		c.JSON(500, gin.H{
			"err": "challengeID: [" + challengeID + "] not found",
		})
		return
	}

	err := challenge.RmChallenge(uid, challengeID)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"remove challenge": "ok",
	})
}
