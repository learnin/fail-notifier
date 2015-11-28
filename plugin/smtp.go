package plugin

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type smtpConfig struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       string
	subject  string
}

type smtpClient struct {
	Host     string
	Port     string
	Username string
	Password string
	client   *smtp.Client
}

type mail struct {
	From    string
	To      string
	Subject string
	Body    string
}

// you should define your plugin struct name of plugin file name.
type Smtp struct {
	config     interface{}
	smtpConfig smtpConfig
}

// you should implement Plugin interface.
func (smtp *Smtp) Notice(s string) {
	smtpClient := smtpClient{
		Host:     smtp.smtpConfig.host,
		Port:     smtp.smtpConfig.port,
		Username: smtp.smtpConfig.username,
		Password: smtp.smtpConfig.password,
	}
	err := smtpClient.connect()
	if err != nil {
		println("SMTP接続に失敗しました。" + err.Error())
		return
	}
	defer func() {
		smtpClient.close()
		smtpClient.quit()
	}()

	smtpClient.sendMail(mail{
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
		smtp.smtpConfig.port = fmt.Sprintf("%v", port)
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

func (c *smtpClient) connect() error {
	client, err := smtp.Dial(fmt.Sprintf("%s:%s", c.Host, c.Port))
	if err != nil {
		return err
	}
	if err := client.Hello("localhost"); err != nil {
		return err
	}
	if c.Username != "" {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(smtp.PlainAuth("", c.Username, c.Password, c.Host)); err != nil {
				return err
			}
		}
	}
	c.client = client
	return nil
}

func (c *smtpClient) close() error {
	return c.client.Close()
}

func (c *smtpClient) quit() error {
	return c.client.Quit()
}

func (c *smtpClient) sendMail(m mail) error {
	if err := c.client.Reset(); err != nil {
		return err
	}
	isNameAddrFrom := false
	var envelopeFrom string
	var fromName string
	r, err := regexp.Compile("^(.*)<(.*)>$")
	if err != nil {
		return err
	}
	if match := r.FindStringSubmatch(m.From); match != nil {
		isNameAddrFrom = true
		fromName = match[1]
		envelopeFrom = match[2]
	} else {
		envelopeFrom = m.From
	}
	c.client.Mail(envelopeFrom)
	if err := c.client.Rcpt(m.To); err != nil {
		return err
	}
	w, err := c.client.Data()
	if err != nil {
		return err
	}
	var headerFrom string
	if isNameAddrFrom {
		name, err := encodeHeader(fromName)
		if err != nil {
			return err
		}
		headerFrom = name + " <" + envelopeFrom + ">"
	} else {
		headerFrom = envelopeFrom
	}

	subject, err := encodeHeader(m.Subject)
	if err != nil {
		return err
	}
	body, err := encodeToJIS(m.Body + "\r\n")
	if err != nil {
		return err
	}
	msg := "From: " + headerFrom + "\r\n" +
		"To: " + m.To + "\r\n" +
		"Subject:" + subject +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=ISO-2022-JP\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" +
		body
	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}
	return w.Close()
}

func encodeToJIS(s string) (string, error) {
	r, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.ISO2022JP.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func encodeHeader(subject string) (string, error) {
	b := make([]byte, 0, utf8.RuneCountInString(subject))
	for _, s := range splitByCharLength(subject, 13) {
		b = append(b, " =?ISO-2022-JP?B?"...)
		s, err := encodeToJIS(s)
		if err != nil {
			return "", err
		}
		b = append(b, base64.StdEncoding.EncodeToString([]byte(s))...)
		b = append(b, "?=\r\n"...)
	}
	return string(b), nil
}

func splitByCharLength(s string, length int) []string {
	result := []string{}
	b := make([]byte, 0, length)
	for i, c := range strings.Split(s, "") {
		b = append(b, c...)
		if i%length == 0 {
			result = append(result, string(b))
			b = make([]byte, 0, length)
		}
	}
	if len(b) > 0 {
		result = append(result, string(b))
	}
	return result
}
