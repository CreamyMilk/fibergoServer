package helper

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg MongoInstance

// Database settings (insert your own database name and connection URI)
const dbName = "fiber_test"
const mongoURI = "mongodb+srv://jotham:yousneekybasterd@cluster0.mtbev.gcp.mongodb.net/testss?retryWrites=true&w=majority"

// Employee struct
type Paper struct {
	Number int    `json:"number"`
	Link   string `json:"link"`
}
type Record struct {
	ID      string  `json:"id,omitempty" bson:"_id,omitempty"`
	PdfName string  `json:"pdfname"`
	NPages  int     `json:"npages"`
	Pages   []Paper `json:"pages"`
}

// Connect configures the MongoDB client and initializes the database connection.
// Source: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup
func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	Mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}
