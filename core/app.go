package core

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// App Structure
type App struct {
	MongoClient 	*mongo.Client
	Err 					error
	DatabaseConnectionContext context.Context
	DatabaseConnectionContextCancelFunc context.CancelFunc
}

// ConnectToMongoDB - establishes connection to the mongo db
func (app *App) ConnectToMongoDB(ping bool) {

	app.MongoClient, app.Err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if app.Err != nil {
		log.Fatalln("Database Connection error:", app.Err)
	}

	app.DatabaseConnectionContext, app.DatabaseConnectionContextCancelFunc =
		context.WithTimeout(context.Background(), 10 * time.Second)

	app.Err = app.MongoClient.Connect(app.DatabaseConnectionContext)
	if ping {
		app.Err = app.MongoClient.Ping(app.DatabaseConnectionContext, readpref.Primary())
	}

	log.Println("Connection established to MongoDB Database")
}

