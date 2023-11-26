package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/posts"
	mw "github.com/twibber/api/app/middleware"
)

func Posts(app fiber.Router) {
	app.Post("/", mw.Auth(true), posts.CreatePost)
	app.Get("/", posts.ListPosts)

	postRouter := app.Group("/:post")
	{
		postRouter.Get("/", posts.GetPost)

		postRouter.Post("/reply", mw.Auth(true), posts.CreateReply)
		postRouter.Post("/repost", mw.Auth(true), posts.CreateRepost)

		postRouter.Post("/like", mw.Auth(true), posts.LikePost)
		postRouter.Delete("/like", mw.Auth(true), posts.UnlikePost)
	}
}
