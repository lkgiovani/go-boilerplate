package delivery

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) ServeOpenAPI(c *fiber.Ctx) error {
	schemaPath := filepath.Join("docs", "schema", "openapi.yaml")
	return c.SendFile(schemaPath)
}

func (h *DocsHandler) ServeRedoc(c *fiber.Ctx) error {
	htmlPath := filepath.Join("resources", "viewers", "redoc.html")
	return c.SendFile(htmlPath)
}

func (h *DocsHandler) ServeSwagger(c *fiber.Ctx) error {
	htmlPath := filepath.Join("resources", "viewers", "swagger.html")
	return c.SendFile(htmlPath)
}

func (h *DocsHandler) ServeScalar(c *fiber.Ctx) error {
	htmlPath := filepath.Join("resources", "viewers", "scalar.html")
	return c.SendFile(htmlPath)
}

func (h *DocsHandler) DocsIndex(c *fiber.Ctx) error {
	htmlPath := filepath.Join("resources", "viewers", "index.html")
	return c.SendFile(htmlPath)
}
