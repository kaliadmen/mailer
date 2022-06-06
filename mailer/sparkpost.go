package mailer

import (
	"github.com/ainsleyclark/go-mail/drivers"
	apimail "github.com/ainsleyclark/go-mail/mail"
	"log"
)

func (m *Mail) useSparkpost(msg Message) error {
	cfg := apimail.Config{
		URL:         m.APIUrl,
		APIKey:      m.APIKey,
		FromAddress: msg.From,
		FromName:    msg.FromName,
	}

	driver, err := drivers.NewSparkPost(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	htmlMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainTextMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	tx := &apimail.Transmission{
		Recipients: []string{msg.To},
		CC:         []string{msg.CC},
		BCC:        []string{msg.BCC},
		Subject:    msg.Subject,
		HTML:       htmlMessage,
		PlainText:  plainTextMessage,
	}

	err = m.addAttachments(msg, tx)
	if err != nil {
		return err
	}

	_, err = driver.Send(tx)
	if err != nil {
		return err
	}

	return nil
}
