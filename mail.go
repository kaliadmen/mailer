package mailer

import (
	"bytes"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

type Mail struct {
	Domain      string
	Templates   string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Jobs        chan Message
	Results     chan Result
	API         string
	APIKey      string
	APIUrl      string
}

type Message struct {
	From        string
	FromName    string
	To          string
	CC          string
	BCC         string
	Subject     string
	Template    string
	Attachments []string
	Data        interface{}
}

type Result struct {
	Success bool
	Error   error
}

func (m *Mail) ListenForMail() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

func (m *Mail) Send(msg Message) error {
	apiLength := len(m.API)
	apiKeyLength := len(m.APIKey)
	apiUrlLength := len(m.APIUrl)
	if apiLength > 0 && apiKeyLength > 0 && apiUrlLength > 0 && m.API != "smtp" {
		err := m.SendUsingAPI(msg)
		if err != nil {
			return err
		}
	}

	return m.SendSMTPMessage(msg)
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	htmlMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainTextMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	timeoutDuration := 10 * time.Second

	mailServer := mail.NewSMTPClient()
	mailServer.Host = m.Host
	mailServer.Port = m.Port
	mailServer.Username = m.Username
	mailServer.Password = m.Password
	mailServer.Encryption = m.getEncryption(m.Encryption)
	mailServer.KeepAlive = false
	mailServer.ConnectTimeout = timeoutDuration
	mailServer.SendTimeout = timeoutDuration

	//connect to mail server
	smtpClient, err := mailServer.Connect()
	if err != nil {
		return err
	}

	//build email message
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).AddCc(msg.CC).AddBcc(msg.BCC).SetSubject(msg.Subject)
	email.SetBody(mail.TextHTML, htmlMessage)
	email.AddAlternative(mail.TextPlain, plainTextMessage)

	//check for attachments
	if len(msg.Attachments) > 0 {
		for _, a := range msg.Attachments {
			email.AddAttachment(a)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)

	tmpl, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tmplBuffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&tmplBuffer, "body", msg.Data); err != nil {
		return "", err
	}

	htmlMsg := tmplBuffer.String()
	htmlMsg, err = m.addInlineCSS(htmlMsg)
	if err != nil {
		return "", err
	}

	return htmlMsg, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.txt.tmpl", m.Templates, msg.Template)

	tmpl, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tmplBuffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&tmplBuffer, "body", msg.Data); err != nil {
		return "", err
	}

	plainText := tmplBuffer.String()

	return plainText, nil
}

func (m *Mail) addInlineCSS(str string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	premail, err := premailer.NewPremailerFromString(str, &opts)
	if err != nil {
		return "", err
	}

	html, err := premail.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	switch encryption {
	case "tls":
		return mail.EncryptionSTARTTLS

	case "ssl":
		return mail.EncryptionSSL

	case "none":
		return mail.EncryptionNone

	default:
		return mail.EncryptionSTARTTLS
	}
}

func (m *Mail) apiSelector(msg Message) error {
	switch m.API {
	case "mailgun":
		err := m.useMailgun(msg)
		if err != nil {
			return err
		}
		return nil

	case "sparkpost":
		err := m.useSparkpost(msg)
		if err != nil {
			return err
		}
		return nil

	case "sendgrid":
		err := m.useSendgrid(msg)
		if err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("unsupported api %s; mailgun, sparkpost, or sendgrid are supported", m.API)
	}
}

func (m *Mail) SendUsingAPI(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	err := m.apiSelector(msg)
	if err != nil {
		return err
	}

	return nil
}
