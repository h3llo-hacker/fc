package routes

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	T "handler/template"
)

func templates(c *gin.Context) {
	templates, err := T.QueryTemplate("all")
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
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
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer file.Close()
	fileContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	fileContent := string(fileContentBytes)
	err = T.InsertTemplate(fileContent, templateName)
	if err != nil {
		c.JSON(500, gin.H{
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
	templateID := c.Request.PostFormValue("id")
	err := T.RemoveTemplate(templateID)
	if err != nil {
		c.JSON(500, gin.H{
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
	templates, err := T.QueryTemplate(templateID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "template enabled",
			"data": templates[0],
		})
	}
}

func templateEnable(c *gin.Context) {
	templateID := c.Param("templateID")
	err := T.EnableTemplate(templateID, true)
	if err != nil {
		c.JSON(500, gin.H{
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
	err := T.EnableTemplate(templateID, false)
	if err != nil {
		c.JSON(500, gin.H{
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

func templateUpdate(c *gin.Context) {
	templateID := c.Param("templateID")
	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer file.Close()
	fileContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	fileContent := string(fileContentBytes)
	err = T.UpdateTemplate(templateID, fileContent)
	if err != nil {
		c.JSON(500, gin.H{
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
