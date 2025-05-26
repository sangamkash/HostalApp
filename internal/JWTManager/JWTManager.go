package JWTManager

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type JWTManager struct {
	signingKey           []byte
	duration             int64 // JWT token exp time in min
	refreshTokenDuration int64 // JWT token exp time in days
}

// NewJWTManager take key and duration in min
func NewJWTManager(jwtKey string, duration int64, refreshTokenDuration int64) *JWTManager {
	return &JWTManager{
		signingKey:           []byte(jwtKey),
		duration:             duration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JWTManager) GenerateToken(userData string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(m.duration) * time.Minute)

	claims := jwt.MapClaims{
		"userData": userData,
		"exp":      expirationTime.Unix(), // exp must be a Unix timestamp
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) GenerateRefreshToken(userData string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(m.refreshTokenDuration*24) * time.Hour)

	claims := jwt.MapClaims{
		"userData": userData,
		"exp":      expirationTime.Unix(), // exp must be a Unix timestamp
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, err
}

// VerifyToken return isExpired,claims,error
func (m *JWTManager) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.signingKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("verifyToken error: %v", err)
	}

	// Extract claims
	claimMap, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Manual expiration check
	if exp, ok := claimMap["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("token has expired")
		}
	} else {
		return nil, fmt.Errorf("expiration claim missing or invalid")
	}
	return claimMap, nil
}

func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	userData, ok := claims["userData"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}
	return m.GenerateToken(userData)
}

//func (m *JWTManager) IsValid(r *http.Request) (bool, jwt.Claims, error) {
//	authHeader := r.Header.Get("Authorization")
//	if authHeader == "" {
//		return false, nil, fmt.Errorf("Missing Authorization header")
//	}
//	// Expecting "Bearer <token>"
//	parts := strings.Split(authHeader, " ")
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return false, nil, fmt.Errorf("Invalid Authorization header")
//	}
//	tokenStr := parts[1]
//
//	// Validate JWT
//	claims, err := m.VerifyToken(tokenStr)
//	if err != nil {
//		return false, nil, err
//	}
//
//	return true, claims, nil
//}

// IsValid this return isExpired ,claims
func (m *JWTManager) IsValid(authHeader string) (jwt.MapClaims, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("invalid Authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid Authorization header")
	}
	tokenStr := parts[1]
	claims, err := m.VerifyToken(tokenStr)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
