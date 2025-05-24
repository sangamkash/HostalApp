package internal

import "github.com/gofiber/fiber/v2"

type APIRoute struct {
	Path    string
	Handler fiber.Handler
}
type IAPIService interface {
	GetFiberRoutes() *[]APIRoute
}
