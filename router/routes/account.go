package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/account"
)

func Account(app fiber.Router) {
	app.Get("/", account.GetAccount)
}