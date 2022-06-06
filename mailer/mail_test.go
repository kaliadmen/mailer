package mailer

import (
	"errors"
	"testing"
)

var msg = Message{
	From:        "test@mitsudo.io",
	FromName:    "Kei",
	To:          "tester@test.com",
	CC:          "tester1@test.com",
	BCC:         "tester2@test.com",
	Subject:     "Test",
	Template:    "test",
	Attachments: []string{"./testdata/mail/test.html.tmpl"},
}

func TestMail_SendSMTPMessage(t *testing.T) {
	err := mailer.SendSMTPMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_SendUsingChannel(t *testing.T) {
	mailer.Jobs <- msg
	res := <-mailer.Results
	if res.Error != nil {
		t.Error(errors.New("failed to send message over channel"))
	}

	msg.To = "invalid_address"
	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error(errors.New("expected error with an invalid To address"))
	}
	msg.To = "tester@test.com"
}

func TestMail_SendUsingAPI(t *testing.T) {
	apimsg := Message{
		To:          "tester@test.com",
		CC:          "tester1@test.com",
		BCC:         "tester2@test.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unsupported"
	mailer.APIKey = "some_key"
	mailer.APIUrl = "https://www.test.com"

	err := mailer.SendUsingAPI(apimsg)
	if err == nil {
		t.Error(errors.New("expected error with an unsupported API"))
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""

}

func TestMail_buildHTMLMessage(t *testing.T) {
	_, err := mailer.buildHTMLMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_buildPlainTextMessage(t *testing.T) {
	_, err := mailer.buildPlainTextMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_Send(t *testing.T) {
	err := mailer.Send(msg)
	if err != nil {
		t.Error(err)
	}

	mailer.API = "unsupported"
	mailer.APIKey = "some_key"
	mailer.APIUrl = "https://www.test.com"

	err = mailer.Send(msg)
	if err == nil {
		t.Error(err)
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_apiSelector(t *testing.T) {
	mailer.API = "unsupported"
	err := mailer.apiSelector(msg)
	if err == nil {
		t.Error(err)
	}
}
