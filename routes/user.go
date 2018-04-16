package routes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/h3llo-hacker/fc/handler/challenge"
	U "github.com/h3llo-hacker/fc/handler/user"
	"github.com/h3llo-hacker/fc/types"
	"github.com/h3llo-hacker/fc/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func users(c *gin.Context) {
	var (
		limit  int
		offset int
		user   U.User
	)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err = strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	items := []string{"UserName", "UserURL", "EmailAddress", "UserID", "Rank", "IsActive", "UserNum"}
	userMap, err := user.QueryUsersRaw(items, limit, offset)
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
	var user types.User

	// Invite Mode
	inviteCode := c.PostForm("invite")
	inviteBy, err := U.GetInvitedBy(inviteCode)
	if err != nil {
		log.Errorf("userCreate GetInvitedBy Error: [%v]", err)
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	user.Invite.InvitedBy = inviteBy

	user.UserName = c.PostForm("username")
	user.Password = utils.Password(c.PostForm("password"))
	user.EmailAddress = c.PostForm("email")
	user.Quota = 1
	user.Register = types.Register_struct{
		IP:     c.PostForm("ip"),
		Region: c.PostForm("ip"),
		System: types.System_struct{
			OS: c.PostForm("os"),
			UA: c.PostForm("ua"),
		},
		Date: time.Now(),
	}
	userID, err := U.AddUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "create user successfully.",
			"data": userID,
		})

		// remove inviter's code
		if user.Invite.InvitedBy != "invite_off" {
			inviter := U.User{
				UserID: user.Invite.InvitedBy,
			}
			err = inviter.RemoveInviteCode(inviteCode)
			if err != nil {
				log.Errorf("RemoveInviteCode error: [%v], UserID: [%v]", err, inviter.UserID)
			}
		}
	}
}

func userUpdate(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
		UserID:  c.Param("userURL"),
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

	query_user_verification, err := user.QueryUser([]string{"Verify"})
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
	if !query_user_verification.Verify.Verification {
		if userValidate.EmailAddress != "" {
			update["EmailAddress"] = userValidate.EmailAddress
		}
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
		UserID:  c.Param("userURL"),
	}
	all := c.Query("all")
	selector := bson.M{
		"Challenges":       0,
		"Login.LastLogins": 0,
		"Password":         0,
		"Register":         0,
	}
	if all == "1" {
		selector = nil
	}
	quser, err := user.QueryUserRaw(selector)
	if err != nil {
		c.JSON(404, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "get user info ok",
			"data": quser,
		})
	}
}

func userChallenges(c *gin.Context) {
	var (
		limit  int
		offset int
	)
	user := U.User{
		UserURL: c.Param("userURL"),
		UserID:  c.Param("userURL"),
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err = strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}

	cType := c.Query("type")
	states := make([]string, 0)
	switch cType {
	case "0": // failed
		states = []string{"failed"}
	case "1": // terminated
		states = []string{"terminated"}
	case "2": // running & created
		states = []string{"running", "created", "creating"}
	case "3": // succeeded
		states = []string{"succeeded"}
	default:
		states = []string{"all"}
	}

	challenges, err := user.QueryUserChallenges(states, limit, offset)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "get user challenges ok",
		"data": challenges,
	})

	// xxxxxxxxxxxxxxx notice!
	if cType != "2" && states[0] != "all" {
		return
	}

	// refreash "created" challenges
	for _, c := range challenges {
		if c.State != "created" && c.State != "creating" {
			continue
		}
		go func(challengeID string) {
			log.Debugf("RefreshChallengeState: [%s]", challengeID)
			err = challenge.RefreshChallengeState(challengeID)
			if err != nil {
				log.Errorf("RefreshChallengeState Error: [%v]", err)
			}
		}(c.ChallengeID)
	}
}

func userFollowers(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
		UserID:  c.Param("userURL"),
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
		UserID:  c.Param("userURL"),
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
	user := U.User{
		UserID: c.Param("userID"),
	}
	targetUser := U.User{
		UserID: c.PostForm("tuid"), //target userID
	}
	action := c.PostForm("action") // follow/unfollow

	if action == "follow" || action == "unfollow" {
		err := user.Follow(action, targetUser)
		if err != nil {
			c.JSON(400, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  action + " user ok",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "action not found!",
		})
	}

}

func userDelete(c *gin.Context) {
	userURL := c.Param("userURL")
	log.Debugf("delete user userURL: [%v]", userURL)
	user := U.User{
		UserURL: userURL,
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
			"code": 1,
			"msg":  "delete user ok",
		})
	}
}

func userLogin(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
		Password:     c.PostForm("password"),
	}
	log.Debugf("email: [%v], pass:[%v]", user.EmailAddress, user.Password)
	userID, success := user.CheckLogin()
	if !success {
		c.JSON(401, gin.H{
			"code": 0,
			"msg":  "login failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "login successfully",
		"data": userID,
	})

	go func() {
		ip := c.PostForm("ip")
		region := utils.IP2Region(ip)
		login := types.Register_struct{
			IP:     ip,
			Region: region,
			System: types.System_struct{
				OS: c.PostForm("os"),
				UA: c.PostForm("ua"),
			},
			Date: time.Now(),
		}
		err := user.UpdateUserLogin(login)
		if err != nil {
			log.Errorf("Update User Login Error: [%v], User: [%v]", err, user.EmailAddress)
		}
	}()
}

func userActive(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
		UserID:  c.Param("userURL"),
	}
	err := user.Active(true)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Active User OK.",
		})
	}
}

func userDeactive(c *gin.Context) {
	user := U.User{
		UserURL: c.Param("userURL"),
		UserID:  c.Param("userURL"),
	}
	err := user.Active(false)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Dective User OK.",
		})
	}
}

func userForgetpasswd(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
	}
	if user.EmailAddress == "" {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "email not found",
		})
		return
	}
	err := user.RestPwd()
	if err != nil {
		log.Errorf("reset passwd error: [%v]", err)
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "forget password.",
	})
}

func userResetpasswd(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
		ResetPwd: types.ResetPwd_struct{
			Code: c.PostForm("code"),
		},
	}
	newPass := c.PostForm("passwd")
	userID, err := user.DoRestPwd(newPass)
	if err != nil {
		log.Errorf("reset password error: [%v]", err)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "reset password error",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "reset password.",
			"data": userID,
		})
	}
}

func userSendVerifyEmail(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
	}
	if user.EmailAddress == "" {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "Email not found",
		})
		return
	}
	err := user.SendVerifyEmail()
	if err != nil {
		// log.Errorf("Send Verify Email error: [%v]", err)
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Verify Email Sent.",
		})
	}
}

func userVerifyEmail(c *gin.Context) {
	user := U.User{
		EmailAddress: c.PostForm("email"),
		Verify: types.Verify_struct{
			Code: c.PostForm("code"),
		},
	}
	err := user.VerifyEmail()
	if err != nil {
		log.Errorf("Verify Email error: [%v]", err)
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "error",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Verify Email.",
		})
	}
}

func userAddInviteCodes(c *gin.Context) {
	user := U.User{
		UserID: c.PostForm("uid"),
	}
	num := c.PostForm("num")
	number, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	}
	err = user.AddInvitecodes(number)
	if err != nil {
		log.Errorf("Verify Email error: [%v]", err)
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "OK",
		})
	}
}
