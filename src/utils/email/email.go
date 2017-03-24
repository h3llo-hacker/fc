package email

import (
	"config"

	log "github.com/Sirupsen/logrus"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendMail() error {
	from := mail.NewEmail("Example User", "mr@xyv6.com")
	subject := "Hello World from the SendGrid Go Library"
	to := mail.NewEmail("Example User", "mr@kfd.me")
	content := mail.NewContent("text/plain", "some text here")
	m := mail.NewV3MailInit(from, subject, to, content)

	if config.Conf.SendGridKey == "" {
		config.LoadConfig()
	}

	request := sendgrid.GetRequest(config.Conf.SendGridKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf("Send Email Error: [%v]", err)
	} else {
		log.Infof("StatusCode: [%v], Body: [%v], Headers: [%v]", response.StatusCode, response.Body, response.Headers)
	}
	return nil
}
