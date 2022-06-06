package mailer

import (
	"github.com/ainsleyclark/go-mail/drivers"
	apimail "github.com/ainsleyclark/go-mail/mail"
	"log"
)

func (m *Mail) useSendgrid(msg Message) error {
	cfg := apimail.Config{
		APIKey:      m.APIKey,
		FromAddress: msg.From,
		FromName:    msg.From,
	}

	driver, err := drivers.NewSendGrid(cfg)
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
