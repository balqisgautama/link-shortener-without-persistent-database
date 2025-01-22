package queries

import (
	"context"
	"time"
	"url-shortener/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLQueries struct {
	collection *mongo.Collection
}

func NewURLQueries(db *mongo.Database) *URLQueries {
	return &URLQueries{
		collection: db.Collection("urls"),
	}
}

func (q *URLQueries) InsertURL(url models.URL) error {
	_, err := q.collection.InsertOne(context.Background(), url)
	return err
}

func (q *URLQueries) FindURL(shortened string) (models.URL, error) {
	var url models.URL
	err := q.collection.FindOne(context.Background(), bson.M{"shortened": shortened}).Decode(&url)
	return url, err
}

func (q *URLQueries) UpdateClickCount(shortened string) error {
	_, err := q.collection.UpdateOne(context.Background(), bson.M{"shortened": shortened}, bson.M{"$inc": bson.M{"clickcount": 1}})
	return err
}

func (q *URLQueries) GetSortedURLs(ascending bool, page int, limit int, isExpired *bool) ([]models.URL, int, error) {
	var urls []models.URL
	sortOrder := 1
	if !ascending {
		sortOrder = -1
	}

	// Calculate the number of documents to skip
	skip := (page - 1) * limit

	// Build the filter
	filter := bson.M{}
	if isExpired != nil {
		if *isExpired {
			filter["expiredat"] = bson.M{"$lt": time.Now().Unix()} // Expired URLs
		} else {
			// Non-expired URLs
			filter["$or"] = []bson.M{
				{"expiredat": bson.M{"$eq": 0}},
				{"expiredat": bson.M{"$gte": time.Now().Unix()}},
			}
		}
	}

	// Get total count of items
	totalCount, err := q.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := q.collection.Find(context.TODO(), filter, options.Find().SetSort(bson.M{"click_count": sortOrder}).SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var url models.URL
		if err := cursor.Decode(&url); err != nil {
			return nil, 0, err
		}
		urls = append(urls, url)
	}
	return urls, int(totalCount), nil
}

func (q *URLQueries) UpdateShortURL(shortURL string, expiry int64, expiredAt int64) error {
	update := bson.M{"$set": bson.M{
		"expiry":    expiry,
		"expiredat": expiredAt,
	}}

	filter := bson.M{
		"shortened": shortURL,
		"$or": []bson.M{
			{"expiredat": bson.M{"$eq": 0}},
			{"expiredat": bson.M{"$gte": time.Now().Unix()}},
		},
	}
	result, err := q.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
