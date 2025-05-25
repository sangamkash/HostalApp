package AuthenticationSystem

import (
	"HostelApp/internal"
	"HostelApp/internal/JWTManager"
	AdminDB "HostelApp/internal/database/Admin"
	"HostelApp/internal/storageData/Admin"
	"github.com/gofiber/fiber/v2"
)

type AuthenticationManager struct {
	dbmanager  *AdminDB.LoginDBManager
	jwtmanager *JWTManager.JWTManager
}

func (m *AuthenticationManager) GetFiberRoutes() *[]internal.APIRoute {
	return &[]internal.APIRoute{
		{"/admin/login", m.login},
		{"/admin/createUser", m.createUser},
	}
}

func NewAuthenticationManager(dbManager *AdminDB.LoginDBManager, jwtManager *JWTManager.JWTManager) *AuthenticationManager {
	instance := &AuthenticationManager{
		dbmanager:  dbManager,
		jwtmanager: jwtManager,
	}
	return instance
}

func (s *AuthenticationManager) login(c *fiber.Ctx) error {
	var user Admin.AdminLogin
	if err := c.BodyParser(&user); err != nil {
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
func (s *AuthenticationManager) createUser(c *fiber.Ctx) error {
	var user Admin.AdminUserDetail
	authHeader := c.Get("Authorization")
	if _, jwtErr := s.jwtmanager.IsValid(authHeader); jwtErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": jwtErr.Error(),
		})
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	if err := s.dbmanager.UserCreate(&user, c.Context()); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"message": "user create",
	}
	return c.JSON(resp)

}
