package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// res := rdb.Get(ctx, "user1")
	// res.Result()
	val, err := rdb.Get(ctx, "user1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("user1", val)

	err = rdb.Set(ctx, "user1", "1", 0).Err()
	if err != nil {
		panic(err)
	}

}
