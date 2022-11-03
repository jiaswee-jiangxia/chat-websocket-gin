package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{Network: "tcp", Addr: "ec2-18-143-165-135.ap-southeast-1.compute.amazonaws.com:30088", Password: "foobared", DB: 5})

	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, "mychannel1")

	// Close the subscription when we are done.
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}
