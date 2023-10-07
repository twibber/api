package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/auth"
)

func Auth(app fiber.Router) {
	app.Post("/login", auth.Login)
	app.Post("/register", auth.Register)
	app.Post("/resend", auth.ResendCode)
	app.Post("/verify", auth.Verify)
}
