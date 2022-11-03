package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func main() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{Network: "tcp", Addr: "ec2-18-143-165-135.ap-southeast-1.compute.amazonaws.com:30088", Password: "foobared", DB: 5})
	err := rdb.Publish(ctx, "mychannel1", "payload").Err()
	if err != nil {
		panic(err)
	}
}
