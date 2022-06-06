package mailer

import (
	"github.com/ainsleyclark/go-mail/drivers"
	apimail "github.com/ainsleyclark/go-mail/mail"
	"io/ioutil"
	"path/filepath"
)

func (m *Mail) useMailgun(msg Message) error {
	cfg := apimail.Config{
		URL:         m.APIUrl,
		APIKey:      m.APIKey,
		FromAddress: msg.From,
		FromName:    msg.FromName,
		Domain:      m.Domain,
	}

	driver, err := drivers.NewMailgun(cfg)
	if err != nil {
		return err
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

	//add attachments
	err = m.addAttachments(msg, tx)
	if err != nil {
		return err
	}

	//send mail
	_, err = driver.Send(tx)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) addAttachments(msg Message, tx *apimail.Transmission) error {
	if len(msg.Attachments) > 0 {
		var attachments []apimail.Attachment

		for _, a := range msg.Attachments {
			var attachment apimail.Attachment
			content, err := ioutil.ReadFile(a)
			if err != nil {
				return err
			}

			filename := filepath.Base(a)
			attachment.Filename = filename
			attachment.Bytes = content

			attachments = append(attachments, attachment)

		}

		tx.Attachments = attachments
	}

	return nil
}
