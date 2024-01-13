package main

import (
	"context"
	"net/http"

	swagger "github.com/Danny-Dasilva/Swagger-Router"
	oasFiber "github.com/Danny-Dasilva/Swagger-Router/support/fiber"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type SwaggerRouter = swagger.Router[oasFiber.HandlerFunc, oasFiber.Route]

const (
	swaggerOpenapiTitle   = "test swagger title"
	swaggerOpenapiVersion = "test swagger version"
)

func main() {
	fiberRouter, oasRouter := setupSwagger()

	err := oasRouter.GenerateAndExposeOpenapi()
	if err != nil {
		panic(err)
	}

	// Serve the Swagger documentation
	fiberRouter.Get(swagger.DefaultJSONDocumentationPath, func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendFile("./docs/swagger.json")
	})

	// Start the Fiber server
	err = fiberRouter.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

func setupSwagger() (*fiber.App, *SwaggerRouter) {
	fiberRouter := fiber.New()

	router, err := swagger.NewRouter(oasFiber.NewRouter(fiberRouter), swagger.Options{
		Context: context.Background(),
		Openapi: &openapi3.T{
			Info: &openapi3.Info{
				Title:   swaggerOpenapiTitle,
				Version: swaggerOpenapiVersion,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	operation := swagger.Operation{}

	_, err = router.AddRawRoute(http.MethodGet, "/hello", okHandler, operation)
	if err != nil {
		panic(err)
	}

	_, err = router.AddRoute(http.MethodPost, "/hello/:value", okHandler, swagger.Definitions{})
	if err != nil {
		panic(err)
	}

	return fiberRouter, router
}

func okHandler(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	return c.SendString("OK")
}
