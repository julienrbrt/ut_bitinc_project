package util

import (
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

//UploadToFTP uploads a given file and a filePath to a FTP server
func UploadToFTP(fileName, filePath string) error {
	//ftp credentials
	ftpServer := os.Getenv("FTP_SERVER")
	ftpUser := os.Getenv("FTP_USERNAME")
	ftpPassword := os.Getenv("FTP_PASSWORD")

	//connect to ftp
	c, err := ftp.Dial(ftpServer, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	//login to ftp
	if err := c.Login(ftpUser, ftpPassword); err != nil {
		return err
	}

	//Upload report to FTP server
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	//change directory to "uploads"
	if err := c.ChangeDir("uploads"); err != nil {
		return err
	}

	//upload file to ftp
	if err := c.Stor(fileName, file); err != nil {
		return err
	}

	//close connection
	if err := c.Quit(); err != nil {
		return err
	}

	return nil
}
