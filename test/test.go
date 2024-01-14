package main

import (
	"context"
	"strings"
	"net/http"

	swagger "github.com/Danny-Dasilva/Swagger-Router"
	oasFiber "github.com/Danny-Dasilva/Swagger-Router/support/fiber"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/rakyll/statik/fs"
	_ "github.com/Danny-Dasilva/CycleTLS" // path to generated statik.go
	"path/filepath"
	"mime"
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
		return c.Status(http.StatusOK).SendFile("./docs/json")
	})

	
    fiberRouter.Use("/*", func(c *fiber.Ctx) error {
        // Get the path of the requested file
        filePath := strings.TrimPrefix(c.Path(), "/")

        // If the path is empty, serve the index file
		if filePath == "" {
			filePath = "index.html"
		}
	
		statikFS, err := fs.New()
		if err != nil {
			panic(err)
		}
	
		// Open the file
		file, err := statikFS.Open("/" + filePath)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}
		defer file.Close()
	
		// Get the MIME type of the file based on its extension
		ext := filepath.Ext(filePath)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			// Default to "application/octet-stream"
			mimeType = "application/octet-stream"
		}
	
		// Set the Content-Type header to the MIME type
		c.Set("Content-Type", mimeType)
	
		// Send the file
		return c.SendStream(file)
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

	_, err = router.AddRoute(http.MethodPost, "/hello3/:value", okHandler, swagger.Definitions{})
	if err != nil {
		panic(err)
	}


	return fiberRouter, router
}

func okHandler(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	return c.SendString("OK")
}
