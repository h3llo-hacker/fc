package routes

import (
	"fmt"
	U "handler/user"
	"time"
	"types"
	"utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func users(c *gin.Context) {
	var user U.User
	items := []string{"UserURL", "UserName", "UserNum",
		"EmailAddress"}
	user.EmailAddress = ""
	userMap, err := user.QueryUserRaw(items)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		log.Errorf("Get All Users Error: [%v]", err)
	} else {
		c.JSON(200, gin.H{
			"Users": userMap,
		})
	}
}

func userCreate(c *gin.Context) {
	var (
		user types.User
		u    U.User
	)

	// Invite Mode
	inviteCode := c.PostForm("invite")
	inviteBy, err := U.GetInvitedBy(inviteCode)
	user.Invite.InvitedBy = inviteBy
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.UserName = c.PostForm("username")
	user.Password = utils.Password(c.PostForm("password"))
	user.EmailAddress = c.PostForm("email")
	user.Quota = 1
	user.Register = types.Register_struct{
		IP:     c.PostForm("ip"),
		Region: c.PostForm("ip"),
		System: types.System_struct{
			OS: c.PostForm("os"),
			UA: c.Request.UserAgent(),
		},
		Date: time.Now(),
	}
	err = U.AddUser(user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		// older user
		u.UserID = user.Invite.InvitedBy
		err = u.RemoveInviteCode(inviteCode)
		if err != nil {
			log.Errorf("RemoveInviteCode error: [%v], UserID: [%v]", err, u.UserID)
		}
		c.JSON(200, gin.H{
			"Add user": "OK",
		})
	}
}

func userUpdate(c *gin.Context) {

}

func userInfo(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"UserName", "UserURL", "EmailAddress", "WebSite", "Intro"}
	quser, err := user.QueryUserRaw(items)
	if err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, quser.(bson.M))
	}
}

func userChallenges(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	var challengeState string

	cType := c.Query("type")
	switch cType {
	case "0": // failed
		challengeState = "failed"
	case "1": // terminated
		challengeState = "terminated"
	case "2": // running
		challengeState = "running"
	default:
		challengeState = "all"
	}

	selector := bson.M{"Challenges": 1}
	quser, err := user.QueryUserWithSelector(selector)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	challenges := quser.Challenges
	if challengeState != "all" {
		var returnChallenges []types.UserChallenge
		for _, c := range challenges {
			if c.State == challengeState {
				returnChallenges = append(returnChallenges, c)
			}
		}
		c.JSON(200, returnChallenges)
	} else {
		c.JSON(200, challenges)
	}
}

func userFollowers(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Followers"}
	quser, err := user.QueryUser(items)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err,
		})
	} else {
		c.JSON(200, gin.H{
			"followers": quser.Followers,
		})
	}
}

func userFollowees(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Followees"}
	quser, err := user.QueryUser(items)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err,
		})
	} else {
		c.JSON(200, gin.H{
			"followers": quser.Following,
		})
	}
}

func userFollow(c *gin.Context) {

}

func userDelete(c *gin.Context) {
	emailAddr, _ := c.GetPostForm("email")
	user := U.User{
		EmailAddress: emailAddr,
	}
	err := user.RmUser()
	if err != nil {
		log.Errorf("Remove User Error: [%v]", err)
		errStr := fmt.Sprintf("Remove User Error: [%v]", err)
		c.JSON(500, gin.H{
			"err": errStr,
		})
	} else {
		c.JSON(200, gin.H{
			"Rm User OK": "OK",
		})
	}
}

func userLogin(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
		Password:     c.PostForm("password"),
	}
	if !user.CheckLogin() {
		c.JSON(401, gin.H{
			"login": "false",
		})
		return
	}
	c.JSON(200, gin.H{
		"login": "true",
	})

	go func() {
		ip := c.ClientIP()
		region := utils.IP2Region(ip)
		login := types.Register_struct{
			IP:     ip,
			Region: region,
			System: types.System_struct{
				OS: "",
				UA: c.Request.UserAgent(),
			},
			Date: time.Now(),
		}
		err := user.UpdateUserLogin(login)
		if err != nil {
			log.Errorf("Update User Login Error: [%v], User: [%v]", err, user.EmailAddress)
		}
	}()
}
