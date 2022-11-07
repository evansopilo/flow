package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

var (
	ErrNoDocument     = errors.New("error no document")
	ErrCreateDocument = errors.New("error create document")
	ErrUpdateDocument = errors.New("error update document")
	ErrDeleteDocument = errors.New("error delete document")
)

func main() {

	app := fiber.New()

	hub, err := eventhub.NewHubFromConnectionString(os.Getenv("COSMOSDB_CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("EVENTHUB_CONNECTION_STRING")))
	if err != nil {
		fmt.Println(err)
		return
	}

	app.Post("/events", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
		defer cancel()

		var event map[string]interface{}

		if err := c.BodyParser(&event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "bad request",
			})
		}

		event_byte, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
		}

		go func() {
			if err := addToEventHub(ctx, hub, event_byte); err != nil {
				log.Println(err)
			}
		}()

		go func() {
			if err := addToCosmosDB(ctx, client, event); err != nil {
				log.Println(err)
			}
		}()

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "event received",
		})
	})

	listenAddr := ":3000"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	app.Listen(listenAddr)
}

func addToEventHub(ctx context.Context, hub *eventhub.Hub, event []byte) error {
	var err = hub.Send(ctx, eventhub.NewEvent(event))
	if err != nil {
		return ErrCreateDocument
	}
	return nil
}

func addToCosmosDB(ctx context.Context, client *mongo.Client, event map[string]interface{}) error {
	var result *mongo.InsertOneResult
	var err error

	coll := client.Database("events").Collection("events")

	result, err = coll.InsertOne(ctx, event)
	if err != nil {
		return ErrCreateDocument
	}

	if result.InsertedID.(string) != event["id"].(string) {
		return ErrCreateDocument
	}
	return nil
}
