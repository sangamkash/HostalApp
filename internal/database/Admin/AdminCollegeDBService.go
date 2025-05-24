package Admin

import (
	"HostelApp/LogColor"
	"HostelApp/LogHelper"
	"HostelApp/internal/storageData/Admin"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"time"
)

type CollegeDBManager struct {
	client            *mongo.Client
	collegeCollection *mongo.Collection
}

func NewCollageDBManager(client *mongo.Client) *CollegeDBManager {
	slog.Info(LogHelper.LogServiceStarting("CollegeDBManager"))
	instance := &CollegeDBManager{
		client: client,
	}
	instance.init()
	slog.Info(LogHelper.LogServiceStarted("CollegeDBManager"))
	return instance
}

func (m *CollegeDBManager) init() {
	m.collegeCollection = m.client.Database("admindb").Collection("collegeConfig")
	err := m.createIndexes()
	if err != nil {
		log.Panicf("!!!panic %v", err)
	}
}

func (m *CollegeDBManager) createIndexes() error {
	// Unique indexes for fields that must be unique
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "collage_unique_name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	// Create all indexes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.collegeCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf(LogColor.Red(fmt.Sprintf("failed to create indexes for loginDB error: %v", err)))
		// Application can still run, but queries will be slower
		// and uniqueness won't be enforced at database level
	}
	return nil
}

func (m *CollegeDBManager) AddCollege(data Admin.CollegeData, ctx context.Context) error {

	filter := bson.M{"collage_unique_name": data.CollageUniqueName}
	err := m.collegeCollection.FindOne(ctx, filter).Err()
	if err == nil {
		return fmt.Errorf("collage already exists")
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		// collage does not exist
		_, insertErr := m.collegeCollection.InsertOne(ctx, data)
		if insertErr != nil {
			return insertErr
		}
		return nil
	}
	return err // some other error
}

func (m *CollegeDBManager) UpdateCollage(data *Admin.CollegeData, ctx context.Context) error {
	result, err := m.collegeCollection.UpdateOne(ctx, bson.M{"collage_unique_name": data.CollageUniqueName}, data)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("college_unique_name not found")
	}
	return nil
}

func (m *CollegeDBManager) DeleteCollage(data *Admin.DelCollegeData, ctx context.Context) error {
	filter := bson.M{"collage_unique_name": data.CollageUniqueName}
	update := bson.M{"$set": bson.M{"mark_as_deleted": true}}

	result, err := m.collegeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with collage_unique_name: %s", data.CollageUniqueName)
	}
	return nil
}

func (m *CollegeDBManager) GetCollege(data *Admin.GetCollegeFilter, ctx context.Context) (string, error) {
	// Initialize empty slice (not nil) to ensure JSON marshals as [] instead of null
	colleges := make([]Admin.CollegeData, 0)

	// Find all non-deleted colleges
	cursor, err := m.collegeCollection.Find(ctx, bson.M{"mark_as_deleted": false})
	if err != nil {
		return "", fmt.Errorf("failed to query colleges: %w", err)
	}

	// Proper defer block with error handling for cursor.Close
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			slog.Warn(fmt.Sprintf("Warning: failed to close MongoDB cursor: %v", cerr))
		}
	}()

	// Decode all results at once (more efficient than one-by-one)
	if err := cursor.All(ctx, &colleges); err != nil {
		return "", fmt.Errorf("failed to decode college data: %w", err)
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(colleges)
	if err != nil {
		return "", fmt.Errorf("failed to marshal college data: %w", err)
	}

	return string(jsonData), nil
}
