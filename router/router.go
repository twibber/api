package router

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	mw "github.com/twibber/api/app/middleware"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/router/routes"
)

func Configure() *fiber.App {
	// configure fiber
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		// TODO: change this later to a custom server header
		ServerHeader: fmt.Sprintf("Twibber"),
		AppName:      "Twibber",
		// error handler
		ErrorHandler: lib.ErrorHandler,
		JSONEncoder: func(v any) ([]byte, error) {
			switch respData := v.(type) {
			case lib.Response:
				body, err := sonic.Marshal(respData)
				if err != nil {
					return nil, err
				}

				var result map[string]any
				if err := sonic.Unmarshal(body, &result); err != nil {
					return nil, err
				}

				var name string
				if respData.ObjectName != "" {
					name = respData.ObjectName
				} else {
					dataValue := reflect.ValueOf(respData.Data)
					if dataValue.IsValid() { // Check if the dataValue is valid before accessing its type
						if dataValue.Kind() == reflect.Ptr {
							dataValue = dataValue.Elem()
						}

						if t := dataValue.Type(); t.Kind() == reflect.Slice {
							name = t.Elem().Name() + "s"
						} else {
							name = t.Name()
						}

						if name == "" || name == "s" {
							name = "data"
						} else {
							name = strings.ToLower(name)
						}
					} else {
						name = "data"
					}
				}

				if data, ok := result["data"]; ok {
					dataValue := reflect.ValueOf(data)
					if dataValue.Kind() == reflect.Slice && dataValue.Len() == 0 {
						result[name] = []any{}
					} else {
						result[name] = data
					}
					delete(result, "data")
				}

				return sonic.Marshal(result)
			}

			return sonic.Marshal(v)
		},
		JSONDecoder: func(data []byte, v any) error {
			return sonic.Unmarshal(data, v)
		},
	})

	// log a successful start
	app.Hooks().OnListen(func(data fiber.ListenData) error {
		log.WithFields(log.Fields{
			"port": data.Port,
			"host": data.Host,
		}).Info("initiated http listener")
		return nil
	})

	// attaches a request ID to help with debugging and supporting users with API errors
	app.Use(requestid.New())

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, lib.Config.Domain)
		},
		AllowCredentials: true,
	}))

	// debug request logger
	app.Use(func(c *fiber.Ctx) error {
		log.WithFields(log.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		}).Debug("request")
		return c.Next()
	})

	// status route
	statusCache := app.Use(cache.New())

	statusCache.All("/", func(c *fiber.Ctx) error {
		var mode = "production"
		if lib.Config.Debug {
			mode = "debug"
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"status": fiber.Map{
				"title":  app.Config().AppName,
				"author": "Petar Markov <petar@nolag.host>",
				"health": "healthy",
				"mode":   mode,
				"time":   time.Now().Unix(),
			},
		})
	})

	routes.Auth(app.Group("/auth"))

	routes.Account(app.Group("/account", mw.Auth(false)))
	routes.Account(app.Group("/user", mw.Auth(false)))

	// Debugging routes
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
