package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/app/controllers/posts"
	"github.com/twibber/api/app/controllers/users"
	mw "github.com/twibber/api/app/middleware"
)

func Users(app fiber.Router) {
	app.Get("/", users.ListUsers)

	userRouter := app.Group("/:user")
	{
		userRouter.Get("/", users.GetUserByUsername)

		userRouter.Get("/posts", posts.GetPostsByUser)

		userRouter.Get("/followers", users.GetFollowersByUsername)
		userRouter.Get("/following", users.GetFollowingByUsername)

		userRouter.Post("/follow", mw.Auth(true), users.FollowUser)
		userRouter.Delete("/follow", mw.Auth(true), users.UnfollowUser)
	}
}
