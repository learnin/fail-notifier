package plugin

import (
	"reflect"

	"github.com/learnin/go-send-mail-iso2022jp"
)

type smtpConfig struct {
	host     string
	port     uint16
	username string
	password string
	from     string
	to       string
	subject  string
}

// you should define your plugin struct name of plugin file name.
type Smtp struct {
	config     interface{}
	smtpConfig smtpConfig
}

// you should implement Plugin interface.
func (smtp *Smtp) Notice(s string) {
	c := mail.SmtpClient{
		Host:     smtp.smtpConfig.host,
		Port:     smtp.smtpConfig.port,
		Username: smtp.smtpConfig.username,
		Password: smtp.smtpConfig.password,
	}
	if err := c.Connect(); err != nil {
		println("SMTP接続に失敗しました。" + err.Error())
		return
	}
	defer func() {
		c.Close()
		c.Quit()
	}()
	c.SendMail(mail.Mail{
		From:    smtp.smtpConfig.from,
		To:      smtp.smtpConfig.to,
		Subject: smtp.smtpConfig.subject,
		Body:    s,
	})
}

func (smtp *Smtp) SetPluginConfig(cfg interface{}) {
	smtp.config = cfg
	c, ok := cfg.(map[string]interface{})
	if !ok {
		// FIXME return error
		return
	}
	// FIXME validation
	if host, ok := c["host"].(string); ok {
		smtp.smtpConfig.host = host
	}
	if port, ok := c["port"].(float64); ok {
		smtp.smtpConfig.port = uint16(port)
	}
	if username, ok := c["username"].(string); ok {
		smtp.smtpConfig.username = username
	}
	if password, ok := c["password"].(string); ok {
		smtp.smtpConfig.password = password
	}
	if from, ok := c["from"].(string); ok {
		smtp.smtpConfig.from = from
	}
	if to, ok := c["to"].(string); ok {
		smtp.smtpConfig.to = to
	}
	if subject, ok := c["subject"].(string); ok {
		smtp.smtpConfig.subject = subject
	}
}

// you should register your plugin struct in init function.
func init() {
	registerPluginType(reflect.TypeOf(Smtp{}))
}
