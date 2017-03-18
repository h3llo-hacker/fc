package routes

import (
	"fmt"
	"handler/challenge"
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
	items := []string{"UserID", "UserURL", "UserName", "UserNum",
		"EmailAddress"}
	user.EmailAddress = ""
	userMap, err := user.QueryUsersRaw(items)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		log.Errorf("Get All Users Error: [%v]", err)
	} else {
		c.JSON(200, gin.H{
			"data": userMap,
			"code": 1,
			"msg":  "get all users ok",
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
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
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
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		// older user
		if user.Invite.InvitedBy != "invite_off" {
			u.UserID = user.Invite.InvitedBy
			err = u.RemoveInviteCode(inviteCode)
			if err != nil {
				log.Errorf("RemoveInviteCode error: [%v], UserID: [%v]", err, u.UserID)
			}
		}
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "create user successfully.",
		})
	}
}

func userUpdate(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}

	userValidate := types.User{
		UserName:     c.PostForm("username"),
		Password:     c.PostForm("password"),
		Intro:        c.PostForm("intro"),
		EmailAddress: c.PostForm("email"),
		WebSite:      c.PostForm("website"),
	}
	err := userValidate.ValidateFormat()
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	// UserName, Password, Intro, EmailAddress, WebSite
	update := make(map[string]interface{}, 0)
	if userValidate.UserName != "" {
		update["UserName"] = userValidate.UserName
	}
	if userValidate.Password != "" {
		encPass := utils.Password(userValidate.Password)
		update["Password"] = encPass
	}
	if userValidate.Intro != "" {
		update["Intro"] = userValidate.Intro
	}
	if userValidate.EmailAddress != "" {
		update["EmailAddress"] = userValidate.EmailAddress
	}
	if userValidate.WebSite != "" {
		update["WebSite"] = userValidate.WebSite
	}
	if len(update) == 0 {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "nothing changed",
		})
		return
	}
	log.Debugf("update format [%v]", update)
	update = bson.M{"$set": update}
	err = user.UpdateUser(update)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "update ok",
		})
	}
}

func userInfo(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"UserName", "UserURL", "EmailAddress", "WebSite", "Intro"}
	quser, err := user.QueryUserRaw(items)
	if err != nil {
		c.JSON(404, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "update ok",
			"data": quser.(bson.M),
		})
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

	challenges, err := user.QueryUserChallenges([]string{challengeState})
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	for _, c := range challenges {
		if c.State != "created" {
			continue
		}
		go func(challengeID string) {
			err = challenge.RefreshChallengeState(challengeID)
			if err != nil {
				log.Errorf("RefreshChallengeState Error: [%v]", err)
			}
		}(c.ChallengeID)
	}

	c.JSON(200, challenges)
}

func userFollowers(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Followers"}
	quser, err := user.QueryUser(items)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "get user followers ok",
			"data": quser.Followers,
		})
	}
}

func userFollowees(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
	}
	items := []string{"Following"}
	quser, err := user.QueryUser(items)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "get user followees ok",
			"data": quser.Following,
		})
	}
}

func userFollow(c *gin.Context) {

}

func userDelete(c *gin.Context) {
	emailAddr := c.PostForm("email")
	log.Debugf("user email: [%v]", emailAddr)
	user := U.User{
		EmailAddress: emailAddr,
	}
	err := user.RmUser()
	if err != nil {
		errStr := fmt.Sprintf("Remove User Error: [%v]", err)
		log.Error(errStr)
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  errStr,
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "rm user ok",
		})
	}
}

func userLogin(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
		Password:     c.PostForm("password"),
	}
	log.Debugf("email: [%v], pass:[%v]", user.EmailAddress, user.Password)
	if !user.CheckLogin() {
		c.JSON(401, gin.H{
			"code": 0,
			"msg":  "login false",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "login true",
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
