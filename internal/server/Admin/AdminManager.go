package Admin

import (
	"HostelApp/internal"
	"HostelApp/internal/JWTManager"
	"HostelApp/internal/database/Admin"
	"HostelApp/internal/server/Admin/AuthenticationSystem"
)

type AdminManager struct {
	auth *AuthenticationSystem.AuthenticationManager
}

func (a AdminManager) GetFiberRoutes() *[]internal.APIRoute {
	return a.auth.GetFiberRoutes()
}

func NewAdminManager(adminDb *Admin.DbManager) *AdminManager {
	jwtManager := JWTManager.NewJWTManager("", 30)
	return &AdminManager{
		auth: AuthenticationSystem.NewAuthenticationManager(adminDb.LoginDB, jwtManager),
	}
}
