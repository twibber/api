package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/account"
)

func Account(app fiber.Router) {
	app.Get("/", account.GetAccount)
	app.Get("/email", account.GetAccountEmail)

	app.Get("/connections", account.ListConnections)
	connection := app.Group("/connection/:connection")
	{
		connection.Get("/", account.GetConnection)
		connection.Patch("/password", account.UpdateConnectionPassword)
	}

	app.Get("/sessions", account.ListSessions)
	session := app.Group("/session/:session")
	{
		session.Get("/", account.GetSession)
		session.Delete("/", account.DeleteSession)
	}

	app.Post("/image/:type", account.UpdateProfileImages)

	app.Patch("/", account.UpdateProfile)
}
