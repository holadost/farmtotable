package util

import (
	"github.com/golang/glog"
	"net/smtp"
)

func SendEmail(fromID string, toID string, passwd string, body string) error {
	// TODO: Make this a nice HTML based email with logo and what not.
	msg := "From: " + fromID + "\n" +
		"To: " + toID + "\n" +
		"Subject: Hello there\n\n" +
		body
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", fromID, passwd, "smtp.gmail.com"),
		fromID, []string{toID}, []byte(msg))
	if err != nil {
		glog.Errorf("Unable to send email due to err: %v", err)
		return err
	}
	glog.V(1).Infof("Successfully sent email")
	return nil
}
