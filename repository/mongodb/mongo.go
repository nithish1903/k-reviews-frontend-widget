package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"k-reviews-frontend-api/entity"
	"log"
	"os"

	"github.com/matryer/resync"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var onceMongo resync.Once
var mongoConn *mongo.Client

type MongoCon struct {
	Connection *mongo.Client
}

// !IMP!! Change error to fatal
func MongoConnect() *MongoCon {
	onceMongo.Do(func() {
		zap.L().Info("Inside mongoconnect function")
		// Set up MongoDB client options
		// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

		// Connect to MongoDB
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			fmt.Println(err)
			zap.L().Fatal("Not able to connect to mongo", zap.Any("err", err))
		}

		// Check the connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			zap.L().Error("Failed to ping mongo", zap.Any("err", err))
		}
		zap.L().Info("Connected to MongoDB!", zap.Any("connection", client))

		mongoConn = client
	})

	return &MongoCon{Connection: mongoConn}
}

func GetDocuments(client *mongo.Client, db string, collection string, filter bson.M) (map[string]interface{}, error) {

	data := make(map[string]interface{})
	coll := client.Database(db).Collection(collection)

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		zap.L().Error("Failed to find documents:", zap.Any("err", err))
		return data, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal("Failed to decode document:", err)
		}

		jsonStr, err := json.Marshal(result)
		if err != nil {
			zap.L().Error("Unable to marshal json", zap.Any("err", err))
		}

		errj := json.Unmarshal(jsonStr, &data)
		if errj != nil {
			zap.L().Error("Unable to unmarshal json", zap.Any("err", errj))
		}

		zap.L().Debug("Actual data in map", zap.Any("data", data))
	}

	return data, nil

}

func InsertOne(client *mongo.Client, db string, collection string, data interface{}) bool {
	coll := client.Database(db).Collection(collection)
	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		zap.L().Error("Could not insert into db", zap.Any("err", err))
		return false
	}
	zap.L().Info("Inserted document with", zap.Any("id:", result.InsertedID))

	return true
}

func UpsertOne(client *mongo.Client, db string, collection string, filter bson.M, data interface{}) bool {
	coll := client.Database(db).Collection(collection)
	opts := options.Replace().SetUpsert(true)

	result, err := coll.ReplaceOne(context.TODO(), filter, data, opts)
	if err != nil {
		zap.L().Error("Could not upsert document into db", zap.Any("err", err))
		return false
	}

	if result.UpsertedCount > 0 {
		zap.L().Info("Inserted new document with", zap.Any("filter:", filter))
	} else {
		zap.L().Info("Updated document with", zap.Any("filter:", filter))
	}

	return true
}
func GetData(client *mongo.Client, db string, collection string, filter bson.M) ([]interface{}, bool) {

	data := make([]interface{}, 0)
	coll := client.Database(db).Collection(collection)

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		zap.L().Error("Failed to find documents:", zap.Any("err", err))
		return data, false
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal("Failed to decode document:", err)
		}

		jsonStr, err := json.Marshal(result)
		if err != nil {
			zap.L().Error("Unable to marshal json", zap.Any("err", err))
		}
		// fmt.Println("data-------->", string(jsonStr))

		var obj interface{}
		errj := json.Unmarshal(jsonStr, &obj)
		if errj != nil {
			zap.L().Error("Unable to unmarshal json", zap.Any("err", errj))
		}
		data = append(data, obj)

		zap.L().Debug("Actual data in map", zap.Any("data", data))
	}

	return data, true

}

func GetAggregateDocuments(client *mongo.Client, db string, collection string, pipeline interface{}) ([]bson.M, error) {

	// Access the specified collection
	coll := client.Database(db).Collection(collection)

	// Perform aggregation
	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode aggregation results
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetReviewsDocumentsOnPagination(client *mongo.Client, db string, collection, votesCollection string, filter bson.M, options *options.FindOptions, limit, skip int64) ([]interface{}, int64, error) {
	coll := client.Database(db).Collection(collection)
	votesColl := client.Database(db).Collection(votesCollection)
	results := make([]interface{}, 0)

	cursor, err := coll.Find(context.Background(), filter, options)

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal("Failed to decode document:", err)
		}
		id := result["_id"].(primitive.ObjectID)

		zap.L().Info("id:", zap.Any("id:", id.Hex()))

		// Count upvotes and downvotes for the current product_review_id
		upvotesCount, err := votesColl.CountDocuments(context.Background(), bson.M{"module_id": id.Hex(), "status": 1})
		if err != nil {
			return nil, 0, err
		}

		downvotesCount, err := votesColl.CountDocuments(context.Background(), bson.M{"module_id": id.Hex(), "status": 0})
		if err != nil {
			return nil, 0, err
		}

		// Append upvotes and downvotes count to the current document
		result["upvotes"] = upvotesCount
		result["downvotes"] = downvotesCount
		jsonStr, err := json.Marshal(result)
		if err != nil {
			zap.L().Error("Unable to marshal json", zap.Any("err", err))
		}

		var obj interface{}
		errj := json.Unmarshal(jsonStr, &obj)
		if errj != nil {
			zap.L().Error("Unable to unmarshal json", zap.Any("err", errj))
		}
		results = append(results, obj)

	}

	// Get total count of documents
	total, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil

}

func GetDocumentsOnPagination(client *mongo.Client, db, collection string, filter bson.M, projection bson.M, options *options.FindOptions, limit, skip int64) ([]interface{}, int64, error) {
	coll := client.Database(db).Collection(collection)
	results := make([]interface{}, 0)

	// Set projection in options
	options.SetProjection(projection)
	cursor, err := coll.Find(context.Background(), filter, options)

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal("Failed to decode document:", err)
		}

		jsonStr, err := json.Marshal(result)
		if err != nil {
			zap.L().Error("Unable to marshal json", zap.Any("err", err))
		}

		var obj interface{}
		errj := json.Unmarshal(jsonStr, &obj)
		if errj != nil {
			zap.L().Error("Unable to unmarshal json", zap.Any("err", errj))
		}
		results = append(results, obj)

	}

	// Get total count of documents
	total, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil

}

// func GetRatingCounts(client *mongo.Client, db string, collection string, filter bson.M) (map[int]int64, error) {
// 	coll := client.Database(db).Collection(collection)

// 	pipeline := []bson.M{
// 		{"$match": filter},
// 		{"$group": bson.M{"_id": "$rating", "count": bson.M{"$sum": 1}}},
// 	}

// 	cursor, err := coll.Aggregate(context.Background(), pipeline)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	ratingCounts := make(map[int]int64)

// 	for cursor.Next(context.Background()) {
// 		var result bson.M
// 		if err := cursor.Decode(&result); err != nil {
// 			return nil, err
// 		}

// 		// Convert _id to int32
// 		rating := int(result["_id"].(int32))
// 		count := result["count"].(int64)

// 		// Convert int32 to int before using it as map key
// 		ratingCounts[int(rating)] = count
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

// 	return ratingCounts, nil
// }

func GetRatingCounts(client *mongo.Client, db string, collection string, pipeline interface{}) (map[int32]int32, error) {
	// Access the specified collection
	coll := client.Database(db).Collection(collection)

	// Perform aggregation
	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode aggregation results
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	// Create a map to store the count of each rating
	ratingCounts := make(map[int32]int32)
	for _, result := range results {
		rating, ok := result["rating"].(int32)
		if !ok {
			return nil, fmt.Errorf("rating field is not of type int32")
		}
		count, ok := result["count"].(int32)
		if !ok {
			return nil, fmt.Errorf("count field is not of type int32")
		}
		// key := fmt.Sprintf("rating_%d", rating)
		ratingCounts[rating] = count
	}

	// Add missing ratings with count 0
	var rating int32
	for rating = 1; rating <= 5; rating++ {
		// key := fmt.Sprintf("rating_%d", rating)
		if _, found := ratingCounts[rating]; !found {
			ratingCounts[rating] = 0
		}
	}

	return ratingCounts, nil
}

func GetProductReviewsQaDocumentsOnPagination(client *mongo.Client, db string, collection, votesCollection string, filter bson.M, options *options.FindOptions, limit, skip int64) ([]interface{}, int64, error) {
	coll := client.Database(db).Collection(collection)
	votesColl := client.Database(db).Collection(votesCollection)
	results := make([]interface{}, 0)

	cursor, err := coll.Find(context.Background(), filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result entity.ProductQuestionAnswer
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		zap.L().Info("result:", zap.Any("result:", result))

		zap.L().Info("result[answer]:", zap.Any("answerArray:", result.Answer))

		for i := range result.Answer {
			ansID := result.Answer[i].ID
			// Count upvotes and downvotes for the current ans_id
			upvotesCount, err := votesColl.CountDocuments(context.Background(), bson.M{"module_id": ansID, "status": 1})
			if err != nil {
				return nil, 0, err
			}
			downvotesCount, err := votesColl.CountDocuments(context.Background(), bson.M{"module_id": ansID, "status": 0})
			if err != nil {
				return nil, 0, err
			}
			zap.L().Info("upvotesCount:", zap.Any("upvotesCount:", upvotesCount))
			zap.L().Info("downvotesCount:", zap.Any("downvotesCount:", upvotesCount))

			result.Answer[i].Upvotes = upvotesCount
			result.Answer[i].Downvotes = downvotesCount
			results = append(results, result)
		}
		// Update the result map with the modified answerArray
	}

	// Get total count of documents
	total, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
