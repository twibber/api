package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/posts"
)

func Posts(app fiber.Router) {
	app.Post("/", posts.CreatePost)
}
