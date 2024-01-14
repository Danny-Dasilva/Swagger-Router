package fiber

import (
	"github.com/davidebianchi/gswagger/apirouter"
	"github.com/gofiber/fiber/v2"
	swaggerFiles "github.com/swaggo/files/v2"
	// "github.com/gofiber/fiber/v2/middleware/filesystem"
	"strings"
	// "net/http"
	"log"
	// "io/fs"
	"path/filepath"
	"mime"
)

type HandlerFunc = fiber.Handler
type Route = fiber.Router

type fiberRouter struct {
	router fiber.Router
}

func NewRouter(router fiber.Router) apirouter.Router[HandlerFunc, Route] {
	return fiberRouter{
		router: router,
	}
}

func (r fiberRouter) AddRoute(method string, path string, handler HandlerFunc) Route {
	return r.router.Add(method, path, handler)
}

func (r fiberRouter) SwaggerHandler(contentType string, blob []byte) HandlerFunc {
	return func(c *fiber.Ctx) error {
		log.Println("Path: " + c.Path())
		if strings.HasPrefix(c.Path(), "/docs/swagger/") {
			// Extract the path after "/swagger/"
			filePath := strings.TrimPrefix(c.Path(), "/docs/swagger/")
			log.Println("filePath: " + filePath)
			// Use the extracted path to open the corresponding embedded file
			file, err := swaggerFiles.FS.Open(filePath)
			if err != nil {
				log.Println("Error opening embedded file:", err, filePath)
				return err
			}
			defer file.Close()
		
			// Get file information
			fileInfo, err := file.Stat()
			if err != nil {
				log.Println("Error getting file information:", err)
				return err
			}
		
			 // Get the MIME type of the file based on its extension
			 ext := filepath.Ext(filePath)
			 mimeType := mime.TypeByExtension(ext)
			 if mimeType == "" {
				 // Default to "application/octet-stream"
				 mimeType = "application/octet-stream"
			 }
 
			 // Set the Content-Type header to the MIME type
			 c.Set("Content-Type", mimeType)
		 
			 // Send the file as a stream
			 return c.SendStream(file, int(fileInfo.Size()))
        }
		c.Set("Content-Type", contentType)
		return c.Send(blob)
	}
}

func (r fiberRouter) TransformPathToOasPath(path string) string {
	return apirouter.TransformPathParamsWithColon(path)
}
