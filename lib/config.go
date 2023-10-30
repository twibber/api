package lib

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
)

type Configuration struct {
	// Application
	Debug  bool   `env:"DEBUG"`
	Domain string `env:"DOMAIN"`
	Port   string `env:"PORT"`

	// URLs
	BaseURL  string `env:"BASE_URL"`
	PanelURL string `env:"PANEL_URL"`

	// Database
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_DATABASE"`

	// Only required if DEBUG is false
	MailHost     string `env:"MAIL_HOST"`
	MailPort     string `env:"MAIL_PORT"`
	MailSecure   bool   `env:"MAIL_SECURE"`
	MailUsername string `env:"MAIL_AUTH_USERNAME"`
	MailPassword string `env:"MAIL_AUTH_PASSWORD"`
	MailSender   string `env:"MAIL_SENDER"`
	MailReply    string `env:"MAIL_REPLY"`

	// Captcha
	CaptchaPublic string `env:"CAPTCHA_PUBLIC"`
	CaptchaSecret string `env:"CAPTCHA_SECRET"`

	// OAuth providers
	GoogleClient string `env:"GOOGLE_CLIENT_ID"`
	GoogleSecret string `env:"GOOGLE_CLIENT_SECRET"`

	// Sentry
	SentryDSN string `env:"SENTRY_DSN"`
}

var Config = &Configuration{}

func LoadConfiguration(config *Configuration) {
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		env := typeField.Tag.Get("env")

		// Support for boolean fields
		if typeField.Type.Kind() == reflect.Bool {
			val.Field(i).SetBool(os.Getenv(env) == "true")
		} else {
			val.Field(i).SetString(os.Getenv(env))
		}
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn(".env file not loaded, resorting to environment variables")
	}

	LoadConfiguration(Config)
}
