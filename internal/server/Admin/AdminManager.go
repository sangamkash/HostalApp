package Admin

import (
	"HostelApp/internal"
	"HostelApp/internal/JWTManager"
	"HostelApp/internal/database/Admin"
	"HostelApp/internal/server/Admin/AuthenticationSystem"
	"HostelApp/internal/server/Admin/CollegeSystem"
)

type AdminManager struct {
	auth       *AuthenticationSystem.AuthenticationManager
	collageMng *CollegeSystem.CollegeManager
}

func (a AdminManager) GetFiberRoutes() *[]internal.APIRoute {
	authRoutes := a.auth.GetFiberRoutes()
	collegeRoutes := a.collageMng.GetFiberRoutes()

	// Combine both slices into a new one
	allRoutes := append(*authRoutes, *collegeRoutes...)
	return &allRoutes
}

func NewAdminManager(adminDb *Admin.DbManager) *AdminManager {
	jwtManager := JWTManager.NewJWTManager("", 30, 30)
	return &AdminManager{
		auth:       AuthenticationSystem.NewAuthenticationManager(adminDb.LoginDB, jwtManager),
		collageMng: CollegeSystem.NewCollegeManager(adminDb.CollegeDB, jwtManager),
	}
}
