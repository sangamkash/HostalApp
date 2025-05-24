package server

import (
	"github.com/gofiber/fiber/v2"

	"HostelApp/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "HostelApp",
			AppName:      "HostelApp",
		}),

		db: database.New(),
	}

	return server
}
