package routes

import (
	"fmt"
	"handler/challenge"
	"handler/template"
	"handler/user"
	"strconv"
	"utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func challenges(c *gin.Context) {
	var (
		limit  int
		offset int
	)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 5
		offset = 0
	}
	offset, err = strconv.Atoi(c.Query("offset"))
	if err != nil {
		limit = 5
		offset = 0
	}
	challenges, err := challenge.AllChallenges(limit, offset)
	if err != nil {
		c.JSON(400, gin.H{
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
	challenge, err := challenge.QueryChallenge(filter, nil)
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
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "user [" + userID + "] not found",
		})
		return
	}

	if !validateTemplate(templateID) {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "template [" + templateID + "] not found",
		})
		return
	}

	t, err := template.QueryTemplate(templateID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	log.Debug(t.Enable)
	if !t.Enable {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "Template [" + t.Name + "] Not Enabled.",
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
		c.JSON(400, gin.H{
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
			challenge.RmChallenge(userID, challengeID)
			challenge.UpdateChallengeState(challengeID, "failed")
		} else {
			challenge.UpdateChallengeState(challengeID, "created")
		}
	}()
	c.JSON(201, gin.H{
		"code": 1,
		"msg":  "challenge created",
		"data": challengeID,
	})
}

func challengeRemove(c *gin.Context) {
	uid := c.PostForm("uid")
	challengeID := c.PostForm("cid")

	if !validateUser(uid) {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "uid: [" + uid + "] not found",
		})
		return
	}
	if !validateChallenge(challengeID) {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "challengeID: [" + challengeID + "] not found",
		})
		return
	}

	err := challenge.RmChallenge(uid, challengeID)
	if err != nil {
		c.JSON(400, gin.H{
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

func challengeValidateFlag(c *gin.Context) {
	challengeID := c.Param("challengeID")
	flag := c.Request.PostFormValue("flag")
	userID, right := challenge.ValidateFlag(flag, challengeID)
	if right {
		go func() {
			err := challenge.RmChallenge(userID, challengeID)
			if err != nil {
				log.Errorf("finish challenge and remove it error: [%v], challengeID: [%v]", err, challengeID)
			}
			err = challenge.UpdateChallengeState(challengeID, "succeeded")
			if err != nil {
				log.Errorf("Update Challenge State Error: [%v], ChallengeID: [%v]", err, challengeID)
			}
		}()
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "flag is right, remove challenge now.",
		})
	} else {
		c.JSON(202, gin.H{
			"code": 0,
			"msg":  "flag is wrong, try again later.",
		})
	}
}
