package internal

import "github.com/gofiber/fiber/v2"

type HTTPMethod int

// Define enum-like constants
const (
	GET HTTPMethod = iota
	POST
	PUT
	PATCH
	DELETE
	HEAD
	OPTIONS
	TRACE
	CONNECT
)

// String() method to convert enum to string
func (m HTTPMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	case OPTIONS:
		return "OPTIONS"
	default:
		return "UNKNOWN"
	}
}

type APIRoute struct {
	Path    string
	Method  HTTPMethod
	Handler fiber.Handler
}
type IAPIService interface {
	GetFiberRoutes() *[]APIRoute
}
