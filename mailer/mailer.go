package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/twibber/api/lib"
	"gopkg.in/gomail.v2"
	"strconv"

	htmlTemplate "html/template"
	textTemplate "text/template"

	log "github.com/sirupsen/logrus"
)

// Global mailer dialer and templates.
var (
	mailer *gomail.Dialer

	textTmpl *textTemplate.Template
	htmlTmpl *htmlTemplate.Template
)

// init sets up the mailer with the appropriate configuration and parses the templates.
func init() {
	// Convert the configured port to an integer.
	port, err := strconv.Atoi(lib.Config.MailPort)
	if err != nil {
		log.WithError(err).Fatal("mail port could not be converted to an integer")
		return
	}

	// Configure the dialer with the mail server settings.
	mailer = gomail.NewDialer(lib.Config.MailHost, port, lib.Config.MailUsername, lib.Config.MailPassword)

	// Set TLS configuration based on whether a secure connection is required.
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: !lib.Config.MailSecure, ServerName: lib.Config.MailHost}

	// Log the mailer configuration.
	log.WithFields(log.Fields{
		"host":   lib.Config.MailHost + ":" + lib.Config.MailPort,
		"user":   lib.Config.MailUsername,
		"sender": lib.Config.MailSender,
		"secure": lib.Config.MailSecure,
	}).Info("initiated mailer")

	// Parse HTML templates.
	htmlTmpl, err = htmlTemplate.ParseGlob("mailer/templates/html/*")
	if err != nil {
		log.WithError(err).Fatal("html templates could not be parsed")
		return
	}

	// Parse text templates.
	textTmpl, err = textTemplate.ParseGlob("mailer/templates/text/*")
	if err != nil {
		log.WithError(err).Fatal("text templates could not be parsed")
		return
	}

	// Debug logging for loaded templates and mailer details.
	log.Info("parsed mailer templates")
	log.WithField("templates", htmlTmpl.Templates()).Debug("html templates")
	log.WithField("templates", textTmpl.Templates()).Debug("text templates")
	log.WithField("mailer", mailer).Debug("mailer")
}

// Send composes and sends an email with the provided subject, file, and data.
func Send(subject string, file string, data any) error {
	// Marshal the data into JSON for unmarshalling into Defaults struct.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the Defaults struct.
	var defaultData Defaults
	if err := json.Unmarshal(jsonData, &defaultData); err != nil {
		return err
	}

	// In production, send the actual email.
	if !lib.Config.Debug {
		var htmlEmail, textEmail bytes.Buffer

		// Execute the HTML template with the provided data.
		if err := htmlTmpl.ExecuteTemplate(&htmlEmail, file+".html", data); err != nil {
			panic(err)
		}

		// Execute the text template with the provided data.
		if err := textTmpl.ExecuteTemplate(&textEmail, file+".txt", data); err != nil {
			panic(err)
		}

		// Compose the email message with both HTML and text parts.
		msg := gomail.NewMessage()
		msg.SetHeader("Subject", subject)
		msg.SetAddressHeader("To", defaultData.Email, defaultData.Name)
		msg.SetBody("text/plain", textEmail.String())
		msg.AddAlternative("text/html", htmlEmail.String())
		msg.SetAddressHeader("From", lib.Config.MailSender, "Twibber")
		msg.SetAddressHeader("Reply-To", lib.Config.MailReply, "Twibber Support")

		// Send the email message.
		if err := mailer.DialAndSend(msg); err != nil {
			log.WithError(err).WithField("msg", msg).Error("an error occurred while sending email")
			return err
		}
	} else {
		// In debug mode, log the email data instead of sending.
		log.WithFields(log.Fields{
			"file": file,
			"data": data,
		}).Debug("mocked email send")
	}

	return nil
}
