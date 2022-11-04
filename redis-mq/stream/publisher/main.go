package main

import (
	"log"
	"math/rand"

	"github.com/go-redis/redis"
)

func main() {
	log.Println("Publisher started")
	redisClient := redis.NewClient(&redis.Options{Network: "tcp", Addr: "ec2-18-143-165-135.ap-southeast-1.compute.amazonaws.com:30088", Password: "foobared", DB: 4})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis", err)
	}
	log.Println("Connected to Redis server")

	for i := 0; i < 5; i++ {
		err = publishTicketReceivedEvent(redisClient)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func publishTicketReceivedEvent(client *redis.Client) error {
	log.Println("Publishing event to Redis")
	err := client.XAdd(&redis.XAddArgs{
		Stream:       "tickets",
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"whatHappened": string("ticket received"),
			"ticketID":     int(rand.Intn(100000000)),
			"ticketData":   string("some ticket data"),
		},
	}).Err()
	return err
}
