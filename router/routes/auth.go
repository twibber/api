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
	app.Post("/email/verify", mw.Auth(false), email.Verify)
	app.Post("/email/resend", email.ResendCode)

	app.All("/logout", auth.Logout)
}
