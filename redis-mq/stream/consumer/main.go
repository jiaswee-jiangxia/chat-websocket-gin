package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func main() {
	log.Println("Consumer started")
	redisClient := redis.NewClient(&redis.Options{Network: "tcp", Addr: "ec2-18-143-165-135.ap-southeast-1.compute.amazonaws.com:30088", Password: "foobared", DB: 4})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Unbale to connect to Redis", err)
	}
	log.Println("Connected to Redis server")

	subject := "tickets"
	consumersGroup := "tickets-consumer-group"
	err = redisClient.XGroupCreate(subject, consumersGroup, "0").Err()
	if err != nil {
		log.Println(err)
	}

	lastID := "0"
	for {
		entries, err := redisClient.XRead(&redis.XReadArgs{
			Streams: []string{subject, lastID},
			Count:   2,
			Block:   0,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(entries[0].Messages); i++ {
			messageID := entries[0].Messages[i].ID
			values := entries[0].Messages[i].Values
			eventDescription := fmt.Sprintf("%v", values["whatHappened"])
			ticketID := fmt.Sprintf("%v", values["ticketID"])
			ticketData := fmt.Sprintf("%v", values["ticketData"])
			if eventDescription == "ticket received" {
				err := handleNewTicket(ticketID, ticketData)
				if err != nil {
					log.Fatal(err)
				}
				redisClient.XAck(subject, consumersGroup, messageID)
			}
			lastID = messageID
		}
	}
}

func handleNewTicket(ticketID string, ticketData string) error {
	log.Printf("Handling new ticket id : %s data %s\n", ticketID, ticketData)
	return nil
}
