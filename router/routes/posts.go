package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/posts"
	mw "github.com/twibber/api/app/middleware"
)

func Posts(app fiber.Router) {
	app.Post("/", mw.Auth(true), posts.CreatePost)
	app.Get("/", posts.ListPosts)
}
