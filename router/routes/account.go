package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/account"
)

func Account(app fiber.Router) {
	app.Get("/", account.GetAccount)
	app.Get("/email", account.GetAccountEmail)
	app.Get("/connections", account.ListConnections)
	app.Get("/sessions", account.ListSessions)

	app.Post("/image/:type", account.UpdateProfileImages)

	app.Patch("/", account.UpdateProfile)
}
