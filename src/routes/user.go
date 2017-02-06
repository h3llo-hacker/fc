package routes

import (
	"handler/user"
	"time"
	"types"
	"utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func users(c *gin.Context) {
	items := []string{"UserURL", "UserName", "UserNum",
		"EmailAddress", ""}
	userMap, err := handler.QueryUsersRaw("", items)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": err.Error(),
		})
		log.Errorf("Get All Users Error: [%v]", err)
		return
	}
	c.JSON(200, gin.H{
		"Users": userMap,
	})
}

func userCreate(c *gin.Context) {
	var user types.User

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
	err := handler.AddUser(user)
	if err != nil {
		c.JSON(500, gin.H{
			"Add user error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Add user": "OK",
	})
	go func() {
		region := utils.IP2Region(user.Register.IP)
		set := bson.M{"Register.Region": region}
		update := bson.M{"$set": set}
		err = handler.UpdateUser(user.EmailAddress, update)
		if err != nil {
			log.Errorf("Update User Region Error: %v", err)
		}
	}()
}

func userUpdate(c *gin.Context) {

}

func userInfo(c *gin.Context) {

}

func userChallenges(c *gin.Context) {

}

func userFollowers(c *gin.Context) {

}

func userFollowees(c *gin.Context) {

}

func userFollow(c *gin.Context) {

}
