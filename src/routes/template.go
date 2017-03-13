package routes

import (
	// "fmt"
	"io/ioutil"
	// "os"
	// "time"
	// "utils"

	// log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	// "gopkg.in/mgo.v2/bson"
	T "handler/template"
	// "types"
)

func templates(c *gin.Context) {
	templates, err := T.QueryTemplate("all")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, templates)
	}
}

func templateCreate(c *gin.Context) {
	templateName := c.Request.PostFormValue("name")
	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	defer file.Close()
	fileContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	fileContent := string(fileContentBytes)
	err = T.InsertTemplate(fileContent, templateName)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"InsertTemplate": "ok",
	})
}

func templateRemove(c *gin.Context) {
	templateID := c.Request.PostFormValue("id")
	err := T.RemoveTemplate(templateID)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"Remove Template": "ok",
		})
	}
}

func templateQuery(c *gin.Context) {
	templateID := c.Param("templateID")
	templates, err := T.QueryTemplate(templateID)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, templates[0])
	}
}
