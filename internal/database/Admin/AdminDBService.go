package Admin

import (
	"HostelApp/LogHelper"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type DbManager struct {
	client    *mongo.Client
	LoginDB   *LoginDBManager
	CollegeDB *CollegeDBManager
}

func NewService(client *mongo.Client) *DbManager {
	slog.Info(LogHelper.LogServiceStarted("MongoDBManager"))
	adminDBManager := &DbManager{
		client:    client,
		LoginDB:   NewLoginDBManager(client),
		CollegeDB: NewCollageDBManager(client),
	}
	slog.Info(LogHelper.LogServiceStarted("MongoDBManager"))
	return adminDBManager
}
