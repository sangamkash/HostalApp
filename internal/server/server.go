package server

import (
	"HostelApp/LogColor"
	"HostelApp/internal"
	"HostelApp/internal/server/Admin"
	"github.com/gofiber/fiber/v2"
	"log/slog"

	"HostelApp/internal/database"
)

type FiberServer struct {
	*fiber.App
	db          *database.DBService
	apiServices []internal.IAPIService
}

func New() *FiberServer {
	app := fiber.New(fiber.Config{
		ServerHeader: "HostelAppServer",
		AppName:      "HostelApp",
	})
	db := database.NewDBService()

	server := &FiberServer{
		App: app,
		db:  db,
	}
	server.registerDefaultFiberRoutes()
	slog.Info(LogColor.Yellow("==FiberServer API List=="))
	adminManager := Admin.NewAdminManager(db.AdminDB)
	server.RegisterFiberRoutes(adminManager)
	slog.Info(LogColor.Green("==FiberServer API List=="))
	return server
}
