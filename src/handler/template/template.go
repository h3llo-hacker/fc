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
		Enable:  false,
	}
	err := db.MongoInsert(C, T)
	if err != nil {
		return err
	}

	return nil
}

func QueryAllTemplates() ([]types.Template, error) {
	var (
		query    = bson.M{}
		selector = bson.M{"ID": 1, "Name": 1, "Enable": 1, "Tags": 1}
	)
	templates, err := db.MongoFind(C, query, selector)
	if err != nil {
		return nil, err
	}
	// 404
	if len(templates) == 0 {
		return nil, fmt.Errorf("Templates Not Found")
	}
	Templates := make([]types.Template, len(templates))
	for i, t := range templates {
		tt := t.(bson.M)["Tags"].([]interface{})
		var tags []string
		for _, ttt := range tt {
			tags = append(tags, ttt.(string))
		}
		Templates[i].ID = t.(bson.M)["ID"].(string)
		Templates[i].Name = t.(bson.M)["Name"].(string)
		Templates[i].Tags = tags
		Templates[i].Enable = t.(bson.M)["Enable"].(bool)
		Templates[i].Content = ""
	}
	return Templates, nil
}

func QueryTemplate(templateID string) (types.Template, error) {
	query := bson.M{"ID": templateID}
	templates, err := db.MongoFind(C, query, nil)
	if err != nil {
		return types.Template{}, err
	}
	// 404
	if len(templates) == 0 {
		return types.Template{}, fmt.Errorf("Template [%v] Not Found", templateID)
	}
	// asign []interface{}
	t := templates[0].(bson.M)["Tags"].([]interface{})
	tags := make([]string, len(t))
	for i, tt := range t {
		tags[i] = tt.(string)
	}
	Template := types.Template{
		ID:      templates[0].(bson.M)["ID"].(string),
		Name:    templates[0].(bson.M)["Name"].(string),
		Enable:  templates[0].(bson.M)["Enable"].(bool),
		Content: templates[0].(bson.M)["Content"].(string),
		Tags:    tags,
	}
	return Template, nil
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
	f := fmt.Sprintf("%v", template.Content)
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
