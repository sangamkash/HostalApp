package Admin

import (
	"HostelApp/LogColor"
	"HostelApp/LogHelper"
	"HostelApp/internal/storageData/Admin"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"time"
)

type LoginDBManager struct {
	client         *mongo.Client
	userCollection *mongo.Collection
}

func NewLoginDBManager(client *mongo.Client) *LoginDBManager {
	slog.Info(LogHelper.LogServiceStarting("LoginDBManager"))
	instance := &LoginDBManager{
		client: client,
	}
	instance.init()
	slog.Info(LogHelper.LogServiceStarted("LoginDBManager"))
	return instance
}
func (m *LoginDBManager) init() {
	m.userCollection = m.client.Database("admindb").Collection("adminUsers")
	err := m.createIndexes()
	m.addDefaultData()
	if err != nil {
		log.Panicf("!!!panic %v", err)
	}
}

func (m *LoginDBManager) createIndexes() error {
	// Unique indexes for fields that must be unique
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	// Create all indexes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.userCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf(LogColor.Red(fmt.Sprintf("failed to create indexes for loginDB error: %v", err)))
		// Application can still run, but queries will be slower
		// and uniqueness won't be enforced at database level
	}
	return nil
}

func (m *LoginDBManager) addDefaultData() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.userCollection.FindOne(ctx, bson.M{"username": "admin"}).Err()
	if err != nil {
		slog.Info(LogColor.Pink(fmt.Sprintf("failed to find admin login database error: %v", err)))
		admin := &Admin.AdminUserDetail{
			Username:     "admin",
			Password:     "password@123",
			Email:        "admin@admin.com",
			ExcessLevel:  Admin.Full,
			RefreshToken: "",
		}
		insertErr := m.UserCreate(admin, ctx)
		if insertErr != nil {
			log.Panicf(LogColor.Red(fmt.Sprintf("failed to insert admin user insert error: %v", insertErr)))
		}
	}
}
func (m *LoginDBManager) IsValidCredentials(credentials *Admin.AdminLogin, ctx context.Context) (*string, error) {
	var result bson.M

	// Find user by username
	err := m.userCollection.FindOne(ctx, bson.M{"username": credentials.Username}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found") // user not found
		}
		return nil, fmt.Errorf("internal error: %v", err) // some other DB error
	}

	// Compare passwords (NOTE: consider using hashed passwords in production)
	err = bcrypt.CompareHashAndPassword([]byte(result["password"].(string)), []byte(credentials.Password))
	if err != nil {
		return nil, fmt.Errorf("password mismatch")
	}
	// Get and return _id as string
	objectID, ok := result["_id"].(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("invalid _id format")
	}
	idStr := objectID.Hex()

	return &idStr, nil // success
}
func (m *LoginDBManager) UpdateRefreshToken(_id string, refreshToken string, ctx context.Context) error {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	//update refresh token for the user
	update := bson.M{
		"$set": bson.M{
			"refresh_token": refreshToken,
		},
	}
	result, err := m.userCollection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (m *LoginDBManager) UserExit(userDetail *Admin.AdminUserDetail, ctx context.Context) (bool, error) {
	// Check for existing email or phone
	filter := bson.M{
		"$or": []bson.M{
			{"email": userDetail.Email},
		},
	}
	var existing bson.M
	err := m.userCollection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		return true, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}
	return false, nil
}

func (m *LoginDBManager) UserCreate(userDetail *Admin.AdminUserDetail, ctx context.Context) error {
	// Hash password (never store plain text passwords)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetail.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("false to hash password error: %v", err)
	}

	// Insert new user
	newUser := bson.M{
		"username":     userDetail.Username,
		"password":     string(hashedPassword),
		"email":        userDetail.Email,
		"created":      time.Now(),
		"excess_level": userDetail.ExcessLevel,
		"refreshToken": userDetail.RefreshToken,
	}

	if _, err = m.userCollection.InsertOne(ctx, newUser); err != nil {
		return fmt.Errorf("false to hash password error: %v", err)
	}
	return nil
}

func (m *LoginDBManager) logoutUser(userDetail *Admin.AdminUserDetail, ctx context.Context) error {
	// Hash password (never store plain text passwords)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetail.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("false to hash password error: %v", err)
	}

	// Insert new user
	newUser := bson.M{
		"username":     userDetail.Username,
		"password":     string(hashedPassword),
		"email":        userDetail.Email,
		"created":      time.Now(),
		"refreshToken": userDetail.RefreshToken,
	}

	if _, err = m.userCollection.InsertOne(ctx, newUser); err != nil {
		return fmt.Errorf("false to hash password error: %v", err)
	}
	return nil
}
