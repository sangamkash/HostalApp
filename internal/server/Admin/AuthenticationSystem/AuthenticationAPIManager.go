package AuthenticationSystem

import (
	"HostelApp/internal"
	"HostelApp/internal/JWTManager"
	"HostelApp/internal/ValidatorSystem"
	AdminDB "HostelApp/internal/database/Admin"
	"HostelApp/internal/storageData/Admin"
	"github.com/gofiber/fiber/v2"
)

type AuthenticationManager struct {
	dbManager  *AdminDB.LoginDBManager
	jwtManager *JWTManager.JWTManager
}

func (m *AuthenticationManager) GetFiberRoutes() *[]internal.APIRoute {
	return &[]internal.APIRoute{
		{"/admin/login", internal.POST, m.login},
		{"/admin/User", internal.POST, m.createUser},
		{"/admin/logout", internal.POST, m.logout},
	}
}

func NewAuthenticationManager(dbManager *AdminDB.LoginDBManager, jwtManager *JWTManager.JWTManager) *AuthenticationManager {
	instance := &AuthenticationManager{
		dbManager:  dbManager,
		jwtManager: jwtManager,
	}
	return instance
}

// @Summary Admin login
// @Description Authenticate admin user
// @Tags admin
// @Accept json
// @Produce json
// @Param credentials body Admin.AdminLogin true "Admin credentials"
// @Success 200 {object} map[string]interface{} "Returns JWT token"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/login [post]
func (s *AuthenticationManager) login(c *fiber.Ctx) error {
	var user Admin.AdminLogin
	//Parsing data to user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	//validating the data
	if err := ValidatorSystem.GetValidator().IsValid(&user); err != nil {
		resp := fiber.Map{
			"message": "failed to validate credentials in validator",
			"error":   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	//checking Credentials in DB
	_id, validErr := s.dbManager.IsValidCredentials(&user, c.Context())
	if validErr != nil {
		resp := fiber.Map{
			"message": "failed to validate credentials in DB",
			"error":   validErr.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	//generating new Refresh token
	refreshToken, refreshJwtErr := s.jwtManager.GenerateRefreshToken("")
	if refreshJwtErr != nil {
		resp := fiber.Map{
			"message": "failed to generate refresh JWT",
			"error":   refreshJwtErr.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	//Updating new Refresh token to DB
	dBJwtWriteErr := s.dbManager.UpdateRefreshToken(*_id, refreshToken, c.Context())
	if dBJwtWriteErr != nil {
		resp := fiber.Map{
			"message": "failed to update UpdateRefreshToken in DB ",
			"error":   dBJwtWriteErr.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	//generating new JWT token
	if token, jwtErr := s.jwtManager.GenerateToken(*_id); jwtErr != nil {
		resp := fiber.Map{
			"message": "failed to generate JWT",
			"error":   jwtErr.Error(),
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

// @Summary Create admin user
// @Description Create a new admin user (requires authentication)
// @Tags admin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param user body Admin.AdminUserDetail true "User details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/User [post]
func (s *AuthenticationManager) createUser(c *fiber.Ctx) error {
	var user Admin.AdminUserDetail
	authHeader := c.Get("Authorization")
	if _, jwtErr := s.jwtManager.IsValid(authHeader); jwtErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to validate credentials at jwt",
			"error":   jwtErr.Error(),
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := ValidatorSystem.GetValidator().IsValid(&user); err != nil {
		resp := fiber.Map{
			"message": "failed to validate credentials at validator",
			"error":   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}
	if err := s.dbManager.UserCreate(&user, c.Context()); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"message": "user create",
	}
	return c.Status(fiber.StatusCreated).JSON(resp)

}

// @Summary Logout admin user
// @Description this is for logout (requires authentication)
// @Tags admin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/logout [post]
func (s *AuthenticationManager) logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	_id := ""
	if jwtMap, jwtErr := s.jwtManager.IsValid(authHeader); jwtErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to validate credentials at jwt",
			"error":   jwtErr.Error(),
		})
	} else {
		_id = jwtMap["_id"].(string)
	}

	//Updating new Refresh token to DB
	if dBJwtWriteErr := s.dbManager.UpdateRefreshToken(_id, "", c.Context()); dBJwtWriteErr != nil {
		resp := fiber.Map{
			"message": "failed to update UpdateRefreshToken in DB ",
			"error":   dBJwtWriteErr.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := fiber.Map{
		"message": "logout successfully",
	}
	return c.JSON(resp)
}
