package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/h3llo-hacker/fc/config"
	"github.com/h3llo-hacker/fc/types"
	"github.com/h3llo-hacker/fc/utils/db"

	"github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
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

func InsertTemplate(content interface{}, templateName string) (string, error) {
	if !uniqTemplate(templateName) {
		return "", fmt.Errorf("Template [%v] Exist", templateName)
	}

	uid, _ := uuid.NewV4()
	templateID := fmt.Sprintf("%v", uid)

	T := types.Template{
		ID:      templateID,
		Name:    templateName,
		Content: content,
		Enable:  false,
	}
	err := db.MongoInsert(C, T)
	if err != nil {
		return "", err
	}

	return templateID, nil
}

func QueryAllTemplates(limit, offset int, tags []string) ([]types.Template, error) {
	var (
		query    = bson.M{"Tags": bson.M{"$all": tags}}
		selector = bson.M{"Content": 0}
	)

	if tags[0] == "" {
		query = bson.M{}
	}

	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return nil, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	result := make([]types.Template, 0)
	err = collection.Find(query).Select(selector).Limit(limit).Skip(offset).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func QueryTemplate(templateID string) (types.Template, error) {
	query := bson.M{"ID": templateID}
	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return types.Template{}, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	result := make([]types.Template, 0)
	err = collection.Find(query).Select(nil).All(&result)
	if err != nil {
		return types.Template{}, err
	}
	// 404
	if len(result) != 1 {
		return types.Template{}, fmt.Errorf("Templates Not Found or multi templates founded.")
	}
	return result[0], nil
}

func RemoveTemplate(templateID string) error {
	query := bson.M{"ID": templateID}
	err := db.MongoRemove(C, query)
	if err != nil {
		return err
	}
	return nil
}

func GenerateComposeFile(templateID string, ENV map[string]string) (string, error) {
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
	template_content := fmt.Sprintf("%v", template.Content)
	for k, v := range ENV {
		template_content = strings.Replace(template_content, k, v, 99)
	}
	_, err = file.WriteString(template_content)
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

func TemplateEnable(templateID string, enable bool) error {
	update := bson.M{"$set": bson.M{"Enable": enable}}
	selector := bson.M{"ID": templateID}
	err := db.MongoUpdate(C, selector, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTemplate(templateID string, update bson.M) error {
	selector := bson.M{"ID": templateID}
	err := db.MongoUpdate(C, selector, update)
	if err != nil {
		return err
	}
	return nil
}
