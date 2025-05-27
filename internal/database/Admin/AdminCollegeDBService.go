package Admin

import (
	"HostelApp/LogColor"
	"HostelApp/LogHelper"
	"HostelApp/internal/storageData/Admin"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
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

func (m *CollegeDBManager) addDefaultData() {

}
func (m *CollegeDBManager) AddCollege(colleges *[]Admin.CollegeData, ctx context.Context) ([]Admin.CollegeNameData, error) {
	var addedColleges []Admin.CollegeNameData

	// Start session
	session, sessionErr := m.collegeCollection.Database().Client().StartSession()
	if sessionErr != nil {
		return nil, fmt.Errorf("failed to start session: %v", sessionErr)
	}
	defer session.EndSession(ctx)

	// Run transaction
	startSessionErr := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// Start the transaction
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %v", err)
		}

		for _, college := range *colleges {
			// Check if CollageUniqueName already exists
			count, err := m.collegeCollection.CountDocuments(sc, bson.M{
				"collage_unique_name": college.CollageUniqueName,
				"mark_as_deleted":     false,
			})
			if err != nil {
				return fmt.Errorf("error checking uniqueness for %s: %v", college.CollageUniqueName, err)
			}
			if count > 0 {
				return fmt.Errorf("college with unique name %s already exists", college.CollageUniqueName)
			}

			// Set default value if not explicitly set
			if !college.MarkAsDeleted {
				college.MarkAsDeleted = false
			}

			// Insert the document
			if _, insertErr := m.collegeCollection.InsertOne(sc, college); insertErr != nil {
				return fmt.Errorf("failed to insert college %s: %v", college.CollageUniqueName, insertErr)
			}

			addedColleges = append(addedColleges, Admin.CollegeNameData{
				CollageUniqueName: college.CollageUniqueName,
			})
		}

		// Commit the transaction
		return session.CommitTransaction(sc)
	})

	if startSessionErr != nil {
		return nil, startSessionErr
	}

	return addedColleges, nil
}

// Helper function to validate college data
func validateCollegeData(college *Admin.CollegeData) error {
	validate := validator.New()
	return validate.Struct(college)
}

func (m *CollegeDBManager) UpdateCollage(college *Admin.CollegeData, ctx context.Context) error {
	result, err := m.collegeCollection.UpdateOne(ctx, bson.M{"collage_unique_name": college.CollageUniqueName}, college)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("college_unique_name not found")
	}
	return nil
}

func (m *CollegeDBManager) DeleteCollage(data *Admin.CollegeNameData, ctx context.Context) error {
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

func (m *CollegeDBManager) FetchCollege(filter *Admin.CollegeFilter, ctx context.Context) ([]Admin.CollegeData, error) {
	// Convert Page and Limit
	page := filter.Page
	if page < 1 {
		page = 1
	}

	limit := filter.Limit
	if limit < 1 || limit > 20 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Build MongoDB filter
	query := bson.M{}
	if filter.PinCode != "" {
		query["pin_code"] = filter.PinCode
	}
	query["mark_as_deleted"] = filter.MarkAsDeleted

	// MongoDB find options
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip)
	cursor, err := m.collegeCollection.Find(ctx, query, opts)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	var colleges []Admin.CollegeData
	if err = cursor.All(ctx, &colleges); err != nil {
		return nil, err
	}
	return colleges, nil
}
