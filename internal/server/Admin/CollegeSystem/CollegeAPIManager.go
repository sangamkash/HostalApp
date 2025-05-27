package CollegeSystem

import (
	"HostelApp/internal"
	"HostelApp/internal/JWTManager"
	"HostelApp/internal/ValidatorSystem"
	AdminDB "HostelApp/internal/database/Admin"
	"HostelApp/internal/storageData/Admin"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CollegeManager struct {
	dbManager  *AdminDB.CollegeDBManager
	jwtManager *JWTManager.JWTManager
}

func NewCollegeManager(dbManager *AdminDB.CollegeDBManager, jwtManager *JWTManager.JWTManager) *CollegeManager {
	instance := &CollegeManager{
		dbManager:  dbManager,
		jwtManager: jwtManager,
	}
	return instance
}

func (m *CollegeManager) GetFiberRoutes() *[]internal.APIRoute {
	return &[]internal.APIRoute{
		{"/admin/college", internal.GET, m.GetCollege},
		{"/admin/college", internal.POST, m.AddCollege},
		{"/admin/college", internal.PATCH, m.AddCollege},
	}
}

// @Summary Get college list
// @Description Fetch filtered list of colleges
// @Tags admin
// @Accept json
// @Produce json
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param pin_code query string false "Pin code"
// @Param mark_as_deleted query boolean false "Include deleted items"
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} []Admin.CollegeData
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/college [get]
func (m *CollegeManager) GetCollege(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if _, jwtErr := m.jwtManager.IsValid(authHeader); jwtErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to validate credentials at jwt",
			"error":   jwtErr.Error(),
		})
	}
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	collFilter := Admin.CollegeFilter{
		Page:    page,
		Limit:   limit,
		PinCode: c.Query("pin_code", ""),
	}

	markAsDeleted, err := strconv.ParseBool(c.Query("mark_as_deleted", "false"))
	if err != nil {
		markAsDeleted = false
	}
	collFilter.MarkAsDeleted = markAsDeleted

	// Optional: Validate the filter struct
	if err := ValidatorSystem.GetValidator().IsValid(&collFilter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to validate filter",
			"error":   err.Error(),
		})
	}

	colleges, err := m.dbManager.FetchCollege(&collFilter, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch colleges",
			"error":   err.Error(),
		})
	}

	return c.JSON(colleges)
}

// @Summary Add college
// @Description Add new list of colleges
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param user body []Admin.CollegeData true "College to be added"
// @Success 200 {object} []Admin.CollegeNameData
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/college [post]
func (m *CollegeManager) AddCollege(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if _, jwtErr := m.jwtManager.IsValid(authHeader); jwtErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to validate credentials at jwt",
			"error":   jwtErr.Error(),
		})
	}
	var colleges *[]Admin.CollegeData
	if err := c.BodyParser(colleges); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse colleges",
			"error":   err.Error(),
		})
	}
	if addedColleges, err := m.dbManager.AddCollege(colleges, c.Context()); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to validate credentials at jwt",
			"error":   err.Error(),
		})
	} else {
		return c.JSON(addedColleges)
	}
}

// @Summary update college list
// @Description to update existing collages make sure the unique college name should be present
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param user body Admin.CollegeData true "College to be added"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/college [patch]
func (m *CollegeManager) UpdateCollege(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if _, jwtErr := m.jwtManager.IsValid(authHeader); jwtErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to validate credentials at jwt",
			"error":   jwtErr.Error(),
		})
	}
	var college *Admin.CollegeData
	if err := c.BodyParser(college); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse colleges",
			"error":   err.Error(),
		})
	}
	if err := m.dbManager.UpdateCollage(college, c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update colleges in database",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "college updated",
	})
}
