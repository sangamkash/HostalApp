package database

import (
	"HostelApp/LogColor"
	"HostelApp/LogHelper"
	"HostelApp/internal/database/Admin"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
}

type service struct {
	db    *mongo.Client
	admin *Admin.DbManager
}

var (
	host = os.Getenv("BLUEPRINT_DB_HOST")
	port = os.Getenv("BLUEPRINT_DB_PORT")
)

func New() Service {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Panicf(LogHelper.LogPanic("fail to connect to mongo" + err.Error()))
	}
	if errPing := client.Ping(ctx, nil); errPing != nil {
		log.Panicf(LogColor.Red("!!Panic!! fail to ping MongoDB error: " + errPing.Error()))
		return nil
	}
	return &service{
		db:    client,
		admin: Admin.NewService(client),
	}
}

func (s *service) Health() map[string]string {
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
