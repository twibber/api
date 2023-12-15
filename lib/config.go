package lib

import (
	"github.com/joho/godotenv"   // Package to load .env files
	"github.com/sirupsen/logrus" // Logrus for structured logging
	"os"                         // Standard library package for OS functionality
	"reflect"                    // Standard library package for runtime reflection
)

// Configuration struct to hold all environment variables.
type Configuration struct {
	// General application settings
	Debug bool   `env:"DEBUG"` // Debug mode toggle
	Port  string `env:"PORT"`  // Application port

	// URLs for various services
	Domain    string `env:"DOMAIN"`     // Domain of the application
	APIURL    string `env:"API_URL"`    // API endpoint URL
	PublicURL string `env:"PUBLIC_URL"` // Public facing URL

	// Pterodactyl game server panel URLs and API keys
	PterodactylURL          string `env:"PTERODACTYL_URL"`
	PterodactylAdminAPIKey  string `env:"PTERODACTYL_ADMIN_API_KEY"`
	PterodactylClientAPIKey string `env:"PTERODACTYL_CLIENT_API_KEY"`

	// Database connection details
	DBHost     string `env:"DB_HOST"`     // Database host
	DBPort     string `env:"DB_PORT"`     // Database port
	DBUsername string `env:"DB_USERNAME"` // Database username
	DBPassword string `env:"DB_PASSWORD"` // Database password
	DBName     string `env:"DB_DATABASE"` // Database name

	// Mail server configurations, required if not in debug mode
	MailHost     string `env:"MAIL_HOST"`          // Mail server host
	MailPort     string `env:"MAIL_PORT"`          // Mail server port
	MailSecure   bool   `env:"MAIL_SECURE"`        // Use secure connection
	MailUsername string `env:"MAIL_AUTH_USERNAME"` // Mail server authentication username
	MailPassword string `env:"MAIL_AUTH_PASSWORD"` // Mail server authentication password
	MailSender   string `env:"MAIL_SENDER"`        // Email sender address
	MailReply    string `env:"MAIL_REPLY"`         // Email reply-to address

	// reCAPTCHA keys
	CaptchaPublic string `env:"CAPTCHA_PUBLIC"` // Public key for reCAPTCHA
	CaptchaSecret string `env:"CAPTCHA_SECRET"` // Secret key for reCAPTCHA

	// OAuth providers' credentials
	GoogleClient string `env:"GOOGLE_CLIENT_ID"`     // Google OAuth Client ID
	GoogleSecret string `env:"GOOGLE_CLIENT_SECRET"` // Google OAuth Secret

	DiscordClient  string `env:"DISCORD_CLIENT_ID"`     // Discord OAuth Client ID
	DiscordSecret  string `env:"DISCORD_CLIENT_SECRET"` // Discord OAuth Secret
	DiscordWebhook string `env:"DISCORD_WEBHOOK_URL"`   // Discord webhook URL

	// Sentry error reporting DSN
	SentryDSN string `env:"SENTRY_DSN"` // Data Source Name for Sentry

	// imgproxy Config
	ImgproxyURL  string `env:"IMGPROXY_URL"`  // imgproxy url
	ImgproxyKey  string `env:"IMGPROXY_KEY"`  // imgproxy key
	ImgproxySalt string `env:"IMGPROXY_SALT"` // imgproxy salt
}

// Config holds the global configuration loaded from environment variables.
var Config = &Configuration{}

// LoadConfiguration populates the Config struct with values from environment variables.
func LoadConfiguration(config *Configuration) {
	val := reflect.ValueOf(config).Elem()

	// Iterates over struct fields and sets them with environment variable values.
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		env := typeField.Tag.Get("env")

		if typeField.Type.Kind() == reflect.Bool {
			// Parses and sets boolean fields
			val.Field(i).SetBool(os.Getenv(env) == "true")
		} else {
			// Sets string fields
			val.Field(i).SetString(os.Getenv(env))
		}
	}
}

// init function tries to load the .env file and falls back to system environment variables.
func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn(".env file not loaded, resorting to environment variables")
	}

	// Loads configuration from environment
	LoadConfiguration(Config)
}
