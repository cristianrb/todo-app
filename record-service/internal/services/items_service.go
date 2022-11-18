package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"record-service/internal/logger"
	"record-service/pkg/models"
	"time"
)

var _ ItemsService = &ItemsServiceImpl{}

const dbTimeout = time.Second * 3

type ItemsService interface {
	Insert(item models.TodoItem) error
	GetAll() ([]*models.TodoItem, error)
}

type ItemsServiceImpl struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

type TodoItemEntry struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func NewItemsService(mongoClient *mongo.Client, dbName, collectionName string) *ItemsServiceImpl {
	return &ItemsServiceImpl{
		mongoClient:    mongoClient,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (s *ItemsServiceImpl) Insert(item models.TodoItem) error {
	collection := s.mongoClient.Database(s.dbName).Collection(s.collectionName)
	_, err := collection.InsertOne(context.TODO(), TodoItemEntry{
		Title:       item.Title,
		Description: item.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		logger.Error("Error while trying to insert item", err)
		return err
	}
	return nil
}

func (s *ItemsServiceImpl) GetAll() ([]*models.TodoItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := s.mongoClient.Database(s.dbName).Collection(s.collectionName)
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		logger.Error("Error finding all items", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*models.TodoItem
	for cursor.Next(ctx) {
		var item TodoItemEntry
		err := cursor.Decode(&item)
		if err != nil {
			logger.Error("Error decoding item into slice", err)
			return nil, err
		}
		items = append(items, &models.TodoItem{
			Title:       item.Title,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return items, nil
}
