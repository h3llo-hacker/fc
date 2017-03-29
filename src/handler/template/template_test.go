package template

import (
	"io/ioutil"
	"testing"
)

func Test_InsertTemplate(t *testing.T) {
	composeFile := "/home/mr/test/docker/test-docker-compose.yml"
	fileContent, err := ioutil.ReadFile(composeFile)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	templateName := "testTemplate"
	err = InsertTemplate(string(fileContent), templateName)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}

func Test_QueryTemplateAll(t *testing.T) {
	templates, err := QueryAllTemplates(0, 0, nil)
	if err != nil {
		t.Error(err)
	} else {
		for _, template := range templates {
			t.Log("ID: " + template.ID)
			t.Log("Name: " + template.Name)
			t.Log(template.Content)
			t.Log("--------")
		}
	}
}

func Test_QueryTemplate(t *testing.T) {
	templates, err := QueryAllTemplates(0, 0, nil)
	if err != nil {
		t.Error(err)
	} else {
		for _, template := range templates {
			qt, err := QueryTemplate(template.ID)
			if err != nil {
				t.Log(err)
			}
			if qt.ID == template.ID {
				t.Log("Test_QueryTemplate OK")
			} else {
				t.Error("fail")
			}
		}
	}
}

func Test_GenerateComposeFile(t *testing.T) {
	templates, err := QueryAllTemplates(0, 0, nil)
	if err != nil {
		t.Error(err)
	}
	filepath, err := GenerateComposeFile(templates[0].ID, "flag{hello-world}")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(filepath)
		t.Log(templates[0].Content)
	}
}

func Test_RemoveTemplate(t *testing.T) {
	templates, err := QueryAllTemplates(0, 0, nil)
	if err != nil {
		t.Error(err)
	} else {
		for _, template := range templates {
			err := RemoveTemplate(template.ID)
			if err != nil {
				t.Log(err)
			}
			_, err = QueryTemplate(template.ID)
			if err == nil {
				t.Error("remove faild")
			}
		}
	}
}
