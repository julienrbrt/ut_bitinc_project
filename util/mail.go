package util

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
)

//InformDriver sends a mail attaching the report to a specific email address
func InformDriver(recipient, attachmentPath, startTime, endTime string) error {
	//mail credentials
	mailServer := os.Getenv("MAIL_SERVER")
	mailAddress := os.Getenv("MAIL_EMAIL")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	//build mail
	e := email.NewEmail()
	e.From = fmt.Sprintf("TX2DB Analysis <%s>", mailAddress)
	e.To = []string{recipient}
	e.Subject = "[TX2DB] You have received a new analysis"
	e.Text = []byte(fmt.Sprintf("Hello,\nYour weekly analysis for the period %s to %s is available.\nHave a great day!\n\nThis email has been automatically generated.", startTime, endTime))
	e.AttachFile(attachmentPath)

	err := e.Send(mailServer, LoginAuth(mailAddress, mailPassword))
	if err != nil {
		return err
	}

	return nil
}

//InformInstructor sends a mail attaching all reports to a specific email address
func InformInstructor(startTime, endTime string) error {
	//mail credentials
	mailServer := os.Getenv("MAIL_SERVER")
	mailAddress := os.Getenv("MAIL_EMAIL")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	//recipient
	instructor := os.Getenv("INSTRUCTOR_EMAIL")

	//build mail
	e := email.NewEmail()
	e.From = fmt.Sprintf("TX2DB Analysis <%s>", mailAddress)
	e.To = []string{instructor}
	e.Subject = "[TX2DB] New weekly driver analysis available"
	e.Text = []byte(fmt.Sprintf("Hello,\nThe weekly driver analysis for the period %s to %s are available from the Bolk FTP Server.\nHave a great day!\n\nThis email has been automatically generated.", startTime, endTime))

	if err := e.Send(mailServer, LoginAuth(mailAddress, mailPassword)); err != nil {
		return err
	}

	return nil
}

//InformSystemAdministratorDriverEmailMissing sends a mail to the system administrator
func InformSystemAdministratorDriverEmailMissing(driverPersonID string) error {
	//mail credentials
	mailServer := os.Getenv("MAIL_SERVER")
	mailAddress := os.Getenv("MAIL_EMAIL")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	//recipient
	administrator := os.Getenv("SYSTEM_ADMINISTATOR_EMAIL")

	//build mail
	e := email.NewEmail()
	e.From = fmt.Sprintf("TX2DB Import/Analysis <%s>", mailAddress)
	e.To = []string{administrator}
	e.Subject = "[TX2DB] A driver mail needs to be added"
	e.Text = []byte(fmt.Sprintf("Hello,\nThe driver (personID: %s) does not have an associated email in the TX2DB database. Please add its email in the 'Driver' table so he/she can receive their weekly report.\nHave a great day!\n\nThis email has been automatically generated.", driverPersonID))

	err := e.Send(mailServer, LoginAuth(mailAddress, mailPassword))
	if err != nil {
		return err
	}

	return nil
}

//InformSystemAdministratorFTPError sends a mail to the system administrator
func InformSystemAdministratorFTPError(filePath string) error {
	//mail credentials
	mailServer := os.Getenv("MAIL_SERVER")
	mailAddress := os.Getenv("MAIL_EMAIL")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	//recipient
	administrator := os.Getenv("SYSTEM_ADMINISTATOR_EMAIL")

	//build mail
	e := email.NewEmail()
	e.From = fmt.Sprintf("TX2DB Analysis <%s>", mailAddress)
	e.To = []string{administrator}
	e.Subject = "[TX2DB] FTP upload failed"
	e.Text = []byte(fmt.Sprintf("Hello,\n Something wrong happen while uploading the weekly report to the FTP. Manual upload is hence necessary. The weekly report can be found in _%s_.\nHave a great day!\n\nThis email has been automatically generated.", filePath))

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
