package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/go-redis/redis/v8"
)

const (
	numDeliveries = 100000000
	batchSize     = 100
)

func main() {
	redisClient := redis.NewClient(&redis.Options{Network: "tcp", Addr: "ec2-18-143-165-135.ap-southeast-1.compute.amazonaws.com:30088", Password: "foobared", DB: 0})
	connection, err := rmq.OpenConnectionWithRedisClient("producer", redisClient, nil)
	if err != nil {
		panic(err)
	}

	things, err := connection.OpenQueue("things")
	if err != nil {
		panic(err)
	}
	foobars, err := connection.OpenQueue("foobars")
	if err != nil {
		panic(err)
	}

	var before time.Time
	for i := 0; i < numDeliveries; i++ {
		delivery := fmt.Sprintf("delivery %s", time.Now().Format(time.RFC850))
		if err := things.Publish(delivery); err != nil {
			log.Printf("failed to publish: %s", err)
		}

		if i%batchSize == 0 {
			duration := time.Now().Sub(before)
			before = time.Now()
			perSecond := time.Second / (duration / batchSize)
			log.Printf("produced %d %s %d", i, delivery, perSecond)
			if err := foobars.Publish("foo"); err != nil {
				log.Printf("failed to publish: %s", err)
			}
		}
	}
}
