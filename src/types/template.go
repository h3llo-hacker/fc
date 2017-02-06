package types

type TemplateService struct {
	ID   string
	Name string
	Spec string
}

type Template struct {
	ID       string
	Name     string
	Services []TemplateService
}

/*
  "ID": "xxx-xxx-xxx",
  "Name": "templateName",
  "Services": [
    {
      "ID": "xxx-xxx-xxx",
      "Name": "xxx-xxx-xxx",
      "Spec": {
        "Image": "xxx-xxx-xxx"
      }
    },
    {
      "ID": "xxx-xxx-xxx",
      "Name": "xxx-xxx-xxx",
      "Spec": {
        "Image": "xxx-xxx-xxx"
      }
    }
  ]
*/
