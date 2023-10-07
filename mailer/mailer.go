package mailer

import (
	"bytes"
	"crypto/tls"
	"github.com/bytedance/sonic"
	"github.com/twibber/api/lib"
	"gopkg.in/gomail.v2"
	"strconv"

	htmlTemplate "html/template"
	textTemplate "text/template"

	log "github.com/sirupsen/logrus"
)

var (
	mailer *gomail.Dialer

	textTmpl *textTemplate.Template
	htmlTmpl *htmlTemplate.Template
)

func init() {
	port, err := strconv.Atoi(lib.Config.MailPort)
	if err != nil {
		log.WithError(err).Fatal("mail port could not be converted to an integer")
		return
	}

	mailer = gomail.NewDialer(lib.Config.MailHost, port, lib.Config.MailUsername, lib.Config.MailPassword)

	if lib.Config.MailSecure {
		mailer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: lib.Config.MailHost}
	} else {
		mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true, ServerName: lib.Config.MailHost}
	}

	log.WithFields(log.Fields{
		"host":   lib.Config.MailHost + ":" + lib.Config.MailPort,
		"user":   lib.Config.MailUsername,
		"sender": lib.Config.MailSender,
		"secure": lib.Config.MailSecure,
	}).Info("initiated mailer")

	htmlTmpl, err = htmlTmpl.ParseGlob("mailer/templates/html/*")
	if err != nil {
		log.WithError(err).Fatal("html templates could not be parsed")
		return
	}

	textTmpl, err = textTemplate.ParseGlob("mailer/templates/text/*")
	if err != nil {
		log.WithError(err).Fatal("text templates could not be parsed")
		return
	}

	log.Info("parsed mailer templates")

	log.WithField("templates", htmlTmpl.Templates()).Debug("html templates")
	log.WithField("templates", textTmpl.Templates()).Debug("text templates")
	log.WithField("mailer", mailer).Debug("mailer")
}

func Send(subject string, file string, data any) error {
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		return err
	}

	var defaultData Defaults
	if err := sonic.Unmarshal(jsonData, &defaultData); err != nil {
		return err
	}

	if !lib.Config.Debug {
		var htmlEmail bytes.Buffer
		if err := htmlTmpl.ExecuteTemplate(&htmlEmail, file+".html", data); err != nil {
			panic(err)
		}

		var textEmail bytes.Buffer
		if err := textTmpl.ExecuteTemplate(&textEmail, file+".txt", data); err != nil {
			panic(err)
		}

		msg := gomail.NewMessage()
		msg.SetHeader("Subject", subject)

		msg.SetAddressHeader("To", defaultData.Email, defaultData.Name)

		msg.SetBody("text/plain", textEmail.String())
		msg.AddAlternative("text/html", htmlEmail.String())

		msg.SetAddressHeader("From", lib.Config.MailSender, "Twibber")
		msg.SetAddressHeader("Reply-To", lib.Config.MailReply, "Twibber Support")

		if err := mailer.DialAndSend(msg); err != nil {
			log.WithError(err).WithField("msg", msg).Error(`an error occurred while sending email`)
		}
	} else {
		log.WithFields(log.Fields{
			"file": file,
			"data": data,
		}).Debug("email copy")
	}

	return nil
}
