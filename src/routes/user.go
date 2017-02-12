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
	users, err := user.QueryUserRaw(items)
	if err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, users[0].(bson.M))
	}
}

func userChallenges(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Challenges"}
	challengeType := c.Query("type")
	users, err := user.QueryUser(items)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	switch challengeType {
	case "0": // Failed
		c.JSON(200, gin.H{
			"challenges": users[0].Challenges.Failed,
		})
	case "1": // Finished
		c.JSON(200, gin.H{
			"challenges": users[0].Challenges.Finished,
		})
	case "2": // InProcess
		c.JSON(200, gin.H{
			"challenges": users[0].Challenges.Finished,
		})
	default:
		c.JSON(200, gin.H{
			"challenges": users[0].Challenges,
		})
	}
}

func userFollowers(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Followers"}
	users, err := user.QueryUser(items)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err,
		})
	} else {
		c.JSON(200, gin.H{
			"followers": users[0].Followers,
		})
	}
}

func userFollowees(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Followees"}
	users, err := user.QueryUser(items)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err,
		})
	} else {
		c.JSON(200, gin.H{
			"followers": users[0].Following,
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
