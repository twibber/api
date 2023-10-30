package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/auth"
	"github.com/twibber/api/app/controllers/auth/email"
	mw "github.com/twibber/api/app/middleware"
)

func Auth(app fiber.Router) {
	app.Post("/email/register", email.Register)
	app.Post("/email/login", email.Login)
	app.Post("/email/verify", email.Verify)
	app.Post("/email/resend", email.ResendCode)

	app.Get("/oauth/:provider", auth.AuthorisationURL)
	app.Get("/oauth/:provider/callback", auth.OAuthCallback)

	app.All("/logout", mw.Auth(false), auth.Logout)
}
