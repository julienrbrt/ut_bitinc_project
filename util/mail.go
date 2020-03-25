package util

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
)

//SendReportMail sends a mail attaching the report to a specific email address
func SendReportMail(attachmentPath, startTime, endTime, personID string) error {
	//get mail credentials
	mailServer := os.Getenv("MAIL_SERVER")
	mailAddress := os.Getenv("MAIL_EMAIL")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	//build mail
	e := email.NewEmail()
	e.From = fmt.Sprintf("TX2DB Analysis <%s>", mailAddress)
	e.To = []string{"bit2020@bolk.nl"}
	e.Subject = fmt.Sprintf("[tx2db] Analysis for driver %s available", personID)
	e.Text = []byte(fmt.Sprintf("Hello,\nA new analysis for the driver %s in the period %s to %s is available.\nHave a great day!", personID, startTime, endTime))
	e.AttachFile(attachmentPath)

	err := e.Send(mailServer, LoginAuth(mailAddress, mailPassword))
	if err != nil {
		return err
	}

	return nil
}

//credit https://github.com/go-gomail/gomail/issues/16#issuecomment-73672398
type loginAuth struct {
	username, password string
}

// LoginAuth returns an Auth that implements the LOGIN authentication
// mechanism as defined in RFC 4616.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}
