package main

import (
	"crypto/tls"
	"flag"
	"log"
	"regexp"
	"strings"

	"github.com/tuxaanand/gomail"
)

const emailRegExp = "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

var fromEmail = flag.String("from", "", "FROM email address")
var smtpHost = flag.String("host", "", "SMTP server host")
var smtpPort = flag.Int("port", 25, "SMTP server port")
var smtpUsername = flag.String("user", "", "SMTP Username")
var smtpPassword = flag.String("password", "", "SMTP Password")
var emailBody = flag.String("email", "This is a test mail", "Email content")
var emailSubject = flag.String("subject", "test mail", "Email Subject")
var smtpTimeout = flag.Int("timeout", 10, "SMTP Timeout")

var toEmail = flag.String("to", "", "TO email addresses separated by ';'")
var ssl = flag.Bool("ssl", false, "SSL required")

func main() {
	flag.Parse()

	validate()
	toAddresses := strings.Split(*toEmail, ";")

	m := gomail.NewMessage()
	m.SetHeader("From", *fromEmail)
	m.SetHeader("To", toAddresses...)
	m.SetHeader("Subject", *emailSubject)
	m.SetBody("text/html", *emailBody)

	d := &gomail.Dialer{
		Host:     *smtpHost,
		Port:     *smtpPort,
		Username: *smtpUsername,
		Password: *smtpPassword,
		SSL:      *ssl,
	}

	d.Timeout = smtpTimeout
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Failed sending mail - %v\n", err)
	}

}

func validate() {
	var errors []string

	if len(strings.TrimSpace(*smtpHost)) < 1 {
		errors = append(errors, "Invalid SMTP Host - "+*smtpHost)
	}

	if ok, _ := regexp.MatchString(emailRegExp, *fromEmail); len(strings.TrimSpace(*fromEmail)) < 1 || !ok {
		errors = append(errors, "Invalid FROM Email - "+*fromEmail)
	}

	toAddresses := strings.Split(*toEmail, ";")

	for _, address := range toAddresses {
		if ok, _ := regexp.MatchString(emailRegExp, address); len(strings.TrimSpace(address)) < 1 || !ok {
			errors = append(errors, "Invalid TO Email - "+address)
		}
	}

	if len(strings.TrimSpace(*emailBody)) < 1 {
		errors = append(errors, "Email body cannot be empty")
	}

	if len(errors) > 1 {
		errStr := strings.Join(errors, "\n")
		log.Fatalf("Please check your inputs\n%v", errStr)
	}
}
