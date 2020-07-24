package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "api"
	collectionName = "accounts"
)

var db *mongo.Database

// InitDatabaseClient creates database connection client.
func InitDatabaseClient() (*mongo.Client, error) {
	// create client for DB connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		return nil, err
	}

	// create connection
	if err = client.Connect(context.TODO()); err != nil {
		return nil, err
	}

	// check DB connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	// get database client connection
	client, err := InitDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}

	// prepare signal nitification channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// watch for notification events and clean up on exit
	go func(client *mongo.Client) {
		for {
			switch <-signalChan {
			case syscall.SIGHUP:
				_ = client.Disconnect(context.TODO())
				os.Exit(int(syscall.SIGHUP))
			case syscall.SIGINT:
				_ = client.Disconnect(context.TODO())
				os.Exit(int(syscall.SIGINT))
			case syscall.SIGTERM:
				_ = client.Disconnect(context.TODO())
				os.Exit(int(syscall.SIGTERM))
			case syscall.SIGQUIT:
				_ = client.Disconnect(context.TODO())
				os.Exit(int(syscall.SIGQUIT))
			}
		}
	}(client)

	// set release mode
	gin.SetMode(gin.ReleaseMode)

	// get database handler
	db = client.Database(databaseName)

	// create a gin router with default middleware
	r := gin.Default()

	// register API endpoints
	r.POST("/account/create", AccountCreate)
	r.POST("/account/authenticate", AccountAuthenticate)

	// listen and serve API daemon
	r.Run(":80")
}
