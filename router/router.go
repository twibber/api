package router

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus" // Logrus - Structured logger for Go

	"github.com/gofiber/fiber/v2"                      // Fiber - Express inspired web framework written in Go
	"github.com/gofiber/fiber/v2/middleware/cache"     // Cache middleware for Fiber
	"github.com/gofiber/fiber/v2/middleware/cors"      // CORS middleware for Fiber
	"github.com/gofiber/fiber/v2/middleware/recover"   // Recover middleware for Fiber to handle panics and keep server running
	"github.com/gofiber/fiber/v2/middleware/requestid" // Middleware to attach a request ID for Fiber

	mw "github.com/twibber/api/app/middleware" // Custom middleware for the Twibber application
	"github.com/twibber/api/lib"               // Library containing project-specific configurations
	"github.com/twibber/api/router/routes"     // Package containing route definitions
)

// Configure sets up the Fiber application with various middleware and routes.
func Configure() *fiber.App {
	// Creating a new Fiber application instance with custom configuration.
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,                   // Disables the default startup message to handle it manually
		ServerHeader:          fmt.Sprintf("Twibber"), // Sets the server header, TODO indicates it should be changed later
		AppName:               "Twibber",              // Sets the name of the application
		ErrorHandler:          lib.ErrorHandler,       // Custom error handler for the application
	})

	// Logging the start of the HTTP listener using the app's hooks.
	app.Hooks().OnListen(func(data fiber.ListenData) error {
		log.WithFields(log.Fields{
			"port": data.Port,
			"host": data.Host,
		}).Info("initiated http listener")
		return nil
	})

	// Middleware to attach a unique request ID to each request, aiding in debugging and support.
	app.Use(requestid.New())

	// Middleware to recover from panics and keep the server running.
	app.Use(recover.New())

	// Configuring CORS with a custom function to allow origins that contain the application's domain.
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// this is only a temporary session to allow testing in production
			log.WithFields(log.Fields{
				"origin": origin,
				"domain": lib.Config.Domain,
			}).Warn("cors")

			return strings.Contains(origin, lib.Config.Domain)
		},
		AllowCredentials: true,
	}))

	// Middleware for logging each request in debug mode.
	app.Use(func(c *fiber.Ctx) error {
		log.WithFields(log.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		}).Debug("request")
		return c.Next()
	})

	// Setting up the cache for the status route.
	statusCache := app.Use(cache.New())

	// Status route to provide application health and debug information.
	statusCache.All("/", func(c *fiber.Ctx) error {
		var mode = "production"
		if lib.Config.Debug {
			mode = "debug"
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"status": fiber.Map{
				"title":  app.Config().AppName,
				"author": "Petar Markov <petar@twibber.xyz>",
				"health": "healthy", // hardcoded for now
				"mode":   mode,
				"time":   time.Now().Unix(),
			},
		})
	})

	// Segregate routes
	routes.Auth(app.Group("/auth"))
	routes.Account(app.Group("/account", mw.Auth(false)))
	routes.Account(app.Group("/user", mw.Auth(false)))

	// Debugging block for printing route information, currently commented out.
	/*
		for _, route := range app.GetRoutes() {
			log.WithFields(log.Fields{
				"name":     route.Name,
				"path":     route.Path,
				"params":   route.Params,
				"handlers": route.Handlers,
				"method":   route.Method,
			}).Debug(route.Path)
		}
	*/

	return app
}
