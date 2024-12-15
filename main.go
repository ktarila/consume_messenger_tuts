package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ktarila/messenger_tuts/database"
	"log"
	"time"

	"github.com/elliotchance/phpserialize"
	"github.com/go-redis/redis/v8"
)

type RedisStream struct {
	Fields struct {
		Message string `json:"message"`
	} `json:"fields"`
	ID string `json:"id"`
}

// Outer structure for the JSON object
type Message struct {
	Body    string        `json:"body"`
	Headers []interface{} `json:"headers"` // Assuming headers is an array of unknown types
}

type Body struct {
	ID           int    `json:"id"`
	FunctionName string `json:"functionName"`
}

func main() {
	// Initialize the database
	database.Init("/home/patrick/work/messenger_tuts/var/data.db")
	defer database.DB.Close() // Close the database connection when the program ends

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	stream := "messages" // Your Symfony Redis queue name

	for {
		// Read messages from the stream
		streams, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, "0"},
			Count:   1,
			Block:   50 * time.Millisecond,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				// No new messages within the block timeout
				continue
			}
			log.Fatalf("Error reading stream: %v", err)
		}

		// Process each message
		for _, msgStream := range streams {
			for _, message := range msgStream.Messages {
				fmt.Printf("\nMessage ID: %s\n", message.ID)
				jsonMessage := map[string]interface{}{
					"id":     message.ID,
					"fields": message.Values,
				}

				// Serialize JSON
				jsonData, err := json.Marshal(jsonMessage)
				if err != nil {
					log.Printf("Error serializing message to JSON: %v", err)
					continue
				}

				var redisStream RedisStream
				if err := json.Unmarshal(jsonData, &redisStream); err != nil {
					fmt.Println("Error unmarshalling Redis stream:", err)
					return
				}

				// Parse the actual message body
				var msg string
				err = phpserialize.Unmarshal([]byte(redisStream.Fields.Message), &msg)
				if err != nil {
					log.Printf("Failed to unmarshal message body: %v", err)
					continue
				}

				var content Message
				err = json.Unmarshal([]byte(msg), &content)
				if err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					continue
				}

				var body Body
				err = json.Unmarshal([]byte(content.Body), &body)
				if err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					continue
				}

				database.UpdateShapeArea(body.ID)

				_, err = rdb.XDel(ctx, stream, message.ID).Result()
				if err != nil {
					log.Fatalf("Error deleting message: %v\n", err)
				}
			}
		}
	}
}
