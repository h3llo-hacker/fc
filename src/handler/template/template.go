package template

import (
	"config"
	"fmt"
	"os"
	"strings"
	"types"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	db "utils/db"
)

const C = "templates"

func uniqTemplate(templateName string) bool {
	query := bson.M{"Name": templateName}
	templates, err := db.MongoFind(C, query, nil)
	if err != nil {
		return false
	}
	if len(templates) == 0 {
		return true
	}
	return false
}

func InsertTemplate(content interface{}, templateName string) error {
	if !uniqTemplate(templateName) {
		return fmt.Errorf("Template [%v] Exist", templateName)
	}

	uid, _ := uuid.NewV4()
	templateID := fmt.Sprintf("%v", uid)

	T := types.Template{
		ID:      templateID,
		Name:    templateName,
		Content: content,
	}
	err := db.MongoInsert(C, T)
	if err != nil {
		return err
	}

	return nil
}

func QueryTemplate(templateID string) (t []types.Template, err error) {
	query := bson.M{"ID": templateID}
	if templateID == "all" {
		query = bson.M{}
	}
	templates, err := db.MongoFind(C, query, nil)
	if err != nil {
		return nil, err
	}
	// 404
	if len(templates) == 0 {
		return nil, fmt.Errorf("Template [%v] Not Found", templateID)
	}
	Templates := make([]types.Template, len(templates))
	for i, t := range templates {
		Templates[i].ID = t.(bson.M)["ID"].(string)
		Templates[i].Name = t.(bson.M)["Name"].(string)
		Templates[i].Content = t.(bson.M)["Content"]
	}
	return Templates, nil
}

func RemoveTemplate(templateID string) error {
	query := bson.M{"ID": templateID}
	err := db.MongoRemove(C, query)
	if err != nil {
		return err
	}
	return nil
}

func GenerateComposeFile(templateID, flag string) (string, error) {
	FilePath := config.Conf.ComposeFilePath
	log.Debugln("Composefile path: " + FilePath)
	templateFilePath := FilePath + "/" + templateID + "_docker-compose.yml"
	// create template file
	file, err := os.Create(templateFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	template, err := QueryTemplate(templateID)
	if err != nil {
		return "", err
	}
	f := fmt.Sprintf("%v", template[0].Content)
	ff := strings.Replace(f, "<FLAG>", flag, 99)
	_, err = file.WriteString(ff)
	if err != nil {
		return "", err
	}
	return templateFilePath, nil
}

func TemplateExist(templateID string) bool {
	query := bson.M{"ID": templateID}
	ts, _ := db.MongoFind(C, query, bson.M{"_id": 1})
	if len(ts) == 0 {
		return false
	}
	return true
}
