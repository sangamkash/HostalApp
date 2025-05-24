package Admin

import (
	AdminDB "HostelApp/internal/database/Admin"
	"HostelApp/internal/server"
	"HostelApp/internal/storageData/Admin"
	"github.com/gofiber/fiber/v2"
)

type AuthenticationManager struct {
	app        *fiber.App
	dbmanager  *AdminDB.LoginDBManager
	jwtmanager *server.JWTManager
}

func NewAuthenticationManager(app *fiber.App, dbManager *AdminDB.LoginDBManager, jwtManager *server.JWTManager) *AuthenticationManager {
	instance := &AuthenticationManager{
		app:        app,
		dbmanager:  dbManager,
		jwtmanager: jwtManager,
	}
	return instance
}

func (m *AuthenticationManager) Init() {
	m.app.Post("/admin/login", m.Login)
}

func (s *AuthenticationManager) Login(c *fiber.Ctx) error {
	var user Admin.AdminLogin
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	if err := s.dbmanager.IsValidCredentials(&user, c.Context()); err != nil {
		resp := fiber.Map{
			"message": "failed to validate credentials",
			"error":   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}
	if token, err := s.jwtmanager.GenerateToken(""); err != nil {
		resp := fiber.Map{
			"message": "failed to validate credentials",
			"error":   err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	} else {
		resp := fiber.Map{
			"message":  "successfully login",
			"jwtToken": token,
		}
		return c.JSON(resp)
	}
}
