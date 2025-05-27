package server

import (
	"HostelApp/LogColor"
	"HostelApp/internal"
	"context"
	"fmt"
	"github.com/gofiber/swagger"
	"log"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/contrib/websocket"
)

func (s *FiberServer) registerDefaultFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	s.App.Get("/websocket", websocket.New(s.websocketHandler))
	s.App.Get("/swagger/*", swagger.New(swagger.Config{
		Title: "Hostel API",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))
	slog.Info(LogColor.Orange("==Check Swagger UI for API documentation=="))
	slog.Info(LogColor.Green("http://127.0.0.1:3000/swagger/index.html"))
}

func (s *FiberServer) RegisterFiberRoutes(apiService internal.IAPIService) {
	apiRoute := apiService.GetFiberRoutes()
	for _, route := range *apiRoute {
		slog.Info(LogColor.Orange(route.Method.String()) + ":" + LogColor.Pink(route.Path))
		switch route.Method {
		case internal.GET:
			s.App.Get(route.Path, route.Handler)
			break
		case internal.POST:
			s.App.Post(route.Path, route.Handler)
			break
		case internal.PUT:
			s.App.Put(route.Path, route.Handler)
			break
		case internal.PATCH:
			s.App.Patch(route.Path, route.Handler)
			break
		case internal.DELETE:
			s.App.Delete(route.Path)
			break
		case internal.HEAD:
			s.App.Head(route.Path, route.Handler)
			break
		case internal.OPTIONS:
			s.App.Options(route.Path, route.Handler)
			break
		case internal.TRACE:
			s.App.Trace(route.Path, route.Handler)
			break
		case internal.CONNECT:
			s.App.Connect(route.Path, route.Handler)
			break

		}
	}
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello everyone hostel server is live",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) websocketHandler(con *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			_, _, err := con.ReadMessage()
			if err != nil {
				cancel()
				log.Println("Receiver Closing", err)
				break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
			if err := con.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
				log.Printf("could not write to socket: %v", err)
				return
			}
			time.Sleep(time.Second * 2)
		}
	}
}
