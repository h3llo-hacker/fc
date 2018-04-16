package email

import (
	"config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"types"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func Send(receiver string) error {
	mail := config.Conf.Mail
	API_USER := mail.API.User
	API_KEY := mail.API.Key
	senderName := mail.Sender.Name
	senderAddr := mail.Sender.Addr

	URL := "http://api.sendcloud.net/apiv2/mail/send"
	data := url.Values{
		"apiUser":            {API_USER},
		"apiKey":             {API_KEY},
		"to":                 {receiver},
		"from":               {senderAddr},
		"fromName":           {senderName},
		"subject":            {"ValidateEmail"},
		"templateInvokeName": {mail.Templates.ValidateEmail},
		"respEmailId":        {"false"},
		"html":               {"test email"},
	}

	resp, err := http.PostForm(URL, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("post: %v \n", string(body))

	return nil
}

type EmailRes struct {
	Result     bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func SendEmail(receiver types.ValidateEmail, template string) error {
	mail := config.Conf.Mail
	API_USER := mail.API.User
	API_KEY := mail.API.Key
	senderName := mail.Sender.Name
	senderAddr := mail.Sender.Addr

	sub := bson.M{"%user%": []string{receiver.UserName}, "%href%": []string{receiver.ClickURL}}

	msg := bson.M{"sub": sub, "to": []string{receiver.EmailAddr}}
	xsmtpapi, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("json marshal msg error: [%v]", err)
		return err
	}

	URL := "http://api.sendcloud.net/apiv2/mail/sendtemplate"
	data := url.Values{
		"apiUser":            {API_USER},
		"apiKey":             {API_KEY},
		"from":               {senderAddr},
		"fromName":           {senderName},
		"templateInvokeName": {template},
		"respEmailId":        {"false"},
		"xsmtpapi":           {string(xsmtpapi)},
	}

	resp, err := http.PostForm(URL, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	str := string(body)
	res := EmailRes{}
	json.Unmarshal([]byte(str), &res)
	if !res.Result {
		err = fmt.Errorf("SendMail Error: [%v]", res.Message)
		log.Errorf("%v", err)
		return err
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("SendMail Error: [%v]", res.Message)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SendVerifyEmail(receiver types.ValidateEmail) error {
	log.Infof("Send VerifyEmail: %v", receiver)
	template := config.Conf.Mail.Templates.ValidateEmail
	return SendEmail(receiver, template)
}

func SendResetPwdEmail(receiver types.ValidateEmail) error {
	log.Infof("Send ResetPwdEmail: %v", receiver)
	template := config.Conf.Mail.Templates.ResetPassword
	return SendEmail(receiver, template)
}
