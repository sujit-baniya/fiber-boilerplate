package config

import (
	"github.com/valyala/bytebufferpool"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Mail struct {
	*mail.SMTPServer
	*mail.SMTPClient
	Host        string `mapstructure:"MAIL_HOST" yaml:"host" env:"MAIL_HOST" env-default:"smtp.mailtrap.io"`
	Username    string `mapstructure:"MAIL_USERNAME" yaml:"username" env:"MAIL_USERNAME" env-default:"821c8fc0bb1e19"`
	Password    string `mapstructure:"MAIL_PASSWORD" yaml:"password" env:"MAIL_PASSWORD" env-default:"24edfcaf91afbc"`
	Encryption  string `mapstructure:"MAIL_ENCRYPTION" yaml:"encryption" env:"MAIL_ENCRYPTION" env-default:"tls"`
	FromAddress string `mapstructure:"MAIL_FROM_ADDRESS" yaml:"from_address" env:"MAIL_FROM_ADDRESS" env-default:"itsursujit@gmail.com"`
	FromName    string `mapstructure:"MAIL_FROM_NAME" yaml:"from_name" env:"MAIL_FROM_NAME" env-default:"Verify-Rest"`
	View        *ViewConfig
	Port        int `mapstructure:"MAIL_PORT" yaml:"port" env:"MAIL_PORT" env-default:"2525"`
}

func (m *Mail) Send(to string, subject string, body string, cc string, from string) error {
	if m.SMTPServer == nil {
		m.SetupMailer()
	}
	//New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom(from). //nolint:wsl
				AddTo(to).
				SetSubject(subject)
	if cc != "" { //nolint:wsl
		email.AddCc(cc)
	}
	email.SetBody(mail.TextHTML, body) //nolint:wsl

	//Call Send and pass the client
	err := email.Send(m.SMTPClient)

	if err != nil {
		return err
	} else {
		log.Println("Email Sent")
	}
	return nil
}

func (m *Mail) PrepareHtml(view string, body fiber.Map) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	// app.Settings.Views.Render
	if err := m.View.Template.TemplateEngine.Render(buf, view, body, "layouts/email"); err != nil {
		// handle err
	}
	return buf.String()
}

func (m *Mail) SetupMailer() {

	var err error
	m.SMTPServer = mail.NewSMTPClient()
	m.SMTPServer.Host = m.Host
	m.SMTPServer.Port = m.Port
	m.SMTPServer.Username = m.Username
	m.SMTPServer.Password = m.Password
	if m.Encryption == "tls" {
		m.SMTPServer.Encryption = mail.EncryptionTLS
	} else {
		m.SMTPServer.Encryption = mail.EncryptionSSL
	}

	//Variable to keep alive connection
	m.SMTPServer.KeepAlive = false

	//Timeout for connect to SMTP Server
	m.SMTPServer.ConnectTimeout = 10 * time.Second

	//Timeout for send the data and wait respond
	m.SMTPServer.SendTimeout = 10 * time.Second
	m.SMTPClient, err = m.SMTPServer.Connect()
	if err != nil {
		log.Print(err)
	}
}
