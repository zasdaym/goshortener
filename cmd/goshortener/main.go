package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/zasdaym/goshortener/link"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	listenAddr := flag.String("addr", "0.0.0.0", "Listen address")
	port := flag.String("port", "8080", "Listen port")
	dbURL := flag.String("dburl", "mongodb://127.0.0.1:27017", "MongoDB URL")
	dbName := flag.String("dbname", "goshortener", "MongoDB database name")
	timeout := flag.Duration("timeout", 30*time.Second, "Server request timeout")

	flag.Parse()

	db, err := connectMongoDB(*dbURL, *dbName)
	if err != nil {
		return err
	}

	linkService := link.NewMongoService(db)
	linkServer := link.NewServer(link.ServerOpts{
		Svc:     linkService,
		Timeout: *timeout,
	})
	addr := *listenAddr + ":" + *port
	return linkServer.Start(addr)
}

func connectMongoDB(dbURL string, dbName string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create new mongo client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to mongo database: %w", err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo database: %w", err)
	}
	return client.Database(dbName), nil
}
