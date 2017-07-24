package email

import (
	"bytes"
	"encoding/json"
	"github.com/hhy5861/logrus"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
)

type (
	Mail struct {
		Content string
		Subject string

		Module  string
		Point   string
		Time    string
		Limit   int64
		Message string
	}

	Message struct {
		Module  string
		Point   string
		Name    string
		Time    int64
		Limit   int64
		Message interface{}
	}
)

func NewMail(message *Message) *Mail {
	msg, _ := json.Marshal(message.Message)

	m := &Mail{
		Module:  message.Module,
		Point:   message.Point,
		Time:    time.Unix(message.Time, 0).Format("2006-01-02 03:04:05"),
		Limit:   message.Limit,
		Message: string(msg),
	}

	return m
}

func (m *Mail) SendEmail(id uint, groupId int) {
	Mtemplate := model.EmailTemplate{}
	res, err := Mtemplate.FindById(id)
	if err != nil {
		ps := logrus.Params{
			"err": err,
		}

		logrus.Warn(ps, err, "get email template error")
		return
	}

	m.Content = res.Content
	m.Subject = res.Subject

	result, err := model.GetGroupIdByEmail(groupId)
	if err != nil {
		ps := logrus.Params{
			"err": err,
		}

		logrus.Warn(ps, err, "get send email user list error")
		return
	}

	config := config.Config.Email

	subject := m.subjectParse()
	content := m.contentParse()

	d := gomail.NewDialer(
		config.Smtp,
		config.Port,
		config.User,
		config.Password)

	s, err := d.Dial()
	if err != nil {
		ps := logrus.Params{}
		ps["param"] = config
		logrus.Warn(ps, err, "email Dial error")
	}

	em := gomail.NewMessage()
	for name, addres := range result {
		em.SetHeader("From", config.From)
		em.SetAddressHeader("To", addres, name)
		em.SetHeader("Subject", subject)
		em.SetBody("text/html", content)

		if err := gomail.Send(s, em); err != nil {
			ps := logrus.Params{}
			ps["param"] = m
			logrus.Warn(ps, err, "send email monitor error")
		}
		em.Reset()
	}
}

func (m *Mail) contentParse() string {
	var content bytes.Buffer
	t := template.New("body template")
	t, _ = t.Parse(m.Content)
	t.Execute(&content, m)
	html := content.String()

	return html
}

func (m *Mail) subjectParse() string {
	var subject bytes.Buffer
	t := template.New("subject template")
	t, _ = t.Parse(m.Subject)
	t.Execute(&subject, m)
	html := subject.String()

	return html
}
