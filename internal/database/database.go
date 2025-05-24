package database

import (
	"HostelApp/LogColor"
	"HostelApp/LogHelper"
	"HostelApp/internal/database/Admin"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBService struct {
	db      *mongo.Client
	AdminDB *Admin.DbManager
}

var (
	username = os.Getenv("BLUEPRINT_DB_USERNAME")
	password = os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
)

func NewDBService() *DBService {
	slog.Info(LogHelper.LogServiceStarted("Database"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	slog.Info(LogColor.Yellow("Connecting to MongoDB url:" + uri))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf(LogHelper.LogPanic("fail to connect to mongo" + err.Error()))
	}
	if errPing := client.Ping(ctx, nil); errPing != nil {
		log.Panicf(LogColor.Red("!!Panic!! fail to ping MongoDB error: " + errPing.Error()))
		return nil
	}
	slog.Info(LogHelper.LogServiceStarted("Database"))
	return &DBService{
		db:      client,
		AdminDB: Admin.NewService(client),
	}
}

func (s *DBService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
