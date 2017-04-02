package routes

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	T "handler/template"
)

func templates(c *gin.Context) {
	var (
		limit  int
		offset int
		tags   string
		err    error
	)
	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err = strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}

	tags = c.Query("tags")

	Tags := strings.Split(tags, ",")

	templates, err := T.QueryAllTemplates(limit, offset, Tags)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Get Templates ok",
			"data": templates,
		})
	}
}

func templateCreate(c *gin.Context) {
	templateName := c.Request.PostFormValue("name")
	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer file.Close()
	fileContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	fileContent := string(fileContentBytes)
	err = T.InsertTemplate(fileContent, templateName)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "Insert Template ok",
	})
}

func templateRemove(c *gin.Context) {
	templateID := c.Param("templateID")
	err := T.RemoveTemplate(templateID)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Remove Template ok",
		})
	}
}

func templateQuery(c *gin.Context) {
	templateID := c.Param("templateID")
	template, err := T.QueryTemplate(templateID)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "get template",
			"data": template,
		})
	}
}

func templateEnable(c *gin.Context) {
	templateID := c.Param("templateID")
	err := T.TemplateEnable(templateID, true)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "template enabled failed",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "template enabled",
		})
	}
}

func templateDisable(c *gin.Context) {
	templateID := c.Param("templateID")
	err := T.TemplateEnable(templateID, false)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "template disabled failed",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "template disabled",
		})
	}
}

func templateUpdate(c *gin.Context) {
	var (
		fileContent  string
		templateTags []string
		updates      = make(map[string]interface{}, 0)
	)

	templateID := c.Param("templateID")

	Level := c.Request.PostFormValue("level")
	Tips := c.Request.PostFormValue("tips")
	Details := c.Request.PostFormValue("details")
	templateName := c.Request.PostFormValue("name")
	tags := c.Request.PostFormValue("tags")
	if tags != "" {
		templateTags = strings.Split(tags, ",")
	}
	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		fileContent = ""
	} else {
		defer file.Close()
		fileContentBytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(400, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			return
		}
		fileContent = string(fileContentBytes)
	}

	if fileContent != "" {
		updates["Content"] = fileContent
	}
	if templateTags != nil {
		updates["Tags"] = templateTags
	}
	if templateName != "" {
		updates["Name"] = templateName
	}
	if Level != "" {
		updates["Level"] = Level
	}
	if Tips != "" {
		updates["Tips"] = Tips
	}
	if Details != "" {
		updates["Details"] = Details
	}
	if len(updates) == 0 {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  "update is empty.",
		})
		return
	}

	update := bson.M{"$set": updates}
	err = T.UpdateTemplate(templateID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "Update Template ok",
	})
}
