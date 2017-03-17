package routes

import (
	"fmt"
	"handler/challenge"
	"handler/user"
	"utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func challenges(c *gin.Context) {
	challenges, err := challenge.AllChallenges()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "get all challenges ok",
			"data": challenges,
		})
	}
}

func challengeInfo(c *gin.Context) {
	cid := c.Param("challengeID")
	filter := bson.M{"ID": cid}
	challenge, err := challenge.QueryChallenge(filter)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "get challenge info ok",
			"data": challenge,
		})
	}
}

func challengeCreate(c *gin.Context) {
	userID := c.Request.PostFormValue("uid")
	templateID := c.Request.PostFormValue("templateID")

	if !validateUser(userID) {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "user [" + userID + "] not found",
		})
		return
	}

	if !validateTemplate(templateID) {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "template [" + templateID + "] not found",
		})
		return
	}

	// check quota
	u := user.User{
		UserID: userID,
	}
	tu, err := u.QueryUser([]string{"Quota"})
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	challenges, err := u.QueryUserChallenges([]string{"running",
		"creating", "created"})
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	if len(challenges) >= tu.Quota {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("Quota Not Enough, You've already create [%v] challenges", tu.Quota),
		})
		return
	}

	challengeID := utils.Guuid()

	go func() {
		_, err := challenge.CreateChallenge(userID,
			templateID, challengeID)
		if err != nil {
			log.Errorf("CreateChallenge Error: [%v]", err)
			challenge.UpdateChallengeState(challengeID, "failed")
		} else {
			challenge.UpdateChallengeState(challengeID, "created")
		}
	}()
	c.JSON(201, gin.H{
		"code": 1,
		"msg":  "challenge created",
		"data": challengeID,
		"id":   challengeID,
	})
}

func challengeRemove(c *gin.Context) {
	uid := c.Request.PostFormValue("uid")
	challengeID := c.Request.PostFormValue("cid")

	if !validateUser(uid) {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "uid: [" + uid + "] not found",
		})
		return
	}
	if !validateChallenge(challengeID) {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "challengeID: [" + challengeID + "] not found",
		})
		return
	}

	err := challenge.RmChallenge(uid, challengeID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "challenge removed",
	})
}
