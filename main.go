package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	leaky_bucket_redis "github.com/alibazlamit/leaky_bucket_redis/leaky_bucket"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // Replace with your Redis server password
		DB:       0,                // Replace with your Redis database index
	})
	client.FlushAll(context.Background())

	now := time.Now()
	// Create a new instance of LeakyBucketRedis
	key := "my_bucket"
	rate := 10.0 // 10 requests per second
	threads := 20
	lb := leaky_bucket_redis.NewLeakyBucket(client, key, rate)

	var wg sync.WaitGroup
	// Allow requests
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(requestNum int) {
			waitTime := lb.Allow(context.Background())
			if waitTime == 0 {
				fmt.Println("Request", requestNum, "is allowed")
			} else {
				time.Sleep(waitTime)
				fmt.Println("Request", requestNum, "is allowed after waiting:", waitTime)
			}
			wg.Done()

		}(i)
	}

	wg.Wait()
	fmt.Printf("it took %v to process %v requests on the rate of %f per second", time.Since(now), threads, rate)
}
