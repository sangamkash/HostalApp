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
	"strings"
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

func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	if content, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		if strings.Contains(string(content), "docker") || strings.Contains(string(content), "kubepods") {
			return true
		}
	}

	return false
}

func NewDBService() *DBService {
	slog.Info(LogHelper.LogServiceStarted("Database"))
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	if !IsRunningInDocker() {
		slog.Info(LogColor.Blue("Not Running in Docker ..."))
		host = "localhost"
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	slog.Info(LogColor.Yellow("Connecting to MongoDB url:" + uri + "\n"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf(LogHelper.LogPanic("fail to connect to mongo" + err.Error()))
	}
	if errPing := client.Ping(ctx, nil); errPing != nil {
		newClient, errFallback := fallBack(ctx)
		if errFallback != nil {
			log.Panicf(LogColor.Red("!!Panic!! fail to ping MongoDB error: " + errFallback.Error()))
			return nil
		} else {
			slog.Info(LogHelper.LogServiceStarted("Database fall back"))
			return &DBService{
				db:      newClient,
				AdminDB: Admin.NewService(client),
			}
		}
	}
	slog.Info(LogHelper.LogServiceStarted("Database"))
	return &DBService{
		db:      client,
		AdminDB: Admin.NewService(client),
	}
}

func fallBack(ctx context.Context) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, "localhost", port)
	slog.Info(LogColor.Yellow("Connecting to MongoDB url:" + uri + "\n"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf(LogHelper.LogPanic("fail to connect to mongo" + err.Error()))
		return nil, err
	}
	if errPing := client.Ping(ctx, nil); errPing != nil {
		return nil, err
	}
	return client, nil
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
