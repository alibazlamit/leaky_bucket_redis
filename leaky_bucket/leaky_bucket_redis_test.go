package leaky_bucket_redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestLeakyBucketRedis_Allow(t *testing.T) {
	// Create a Redis client for testing
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create a new LeakyBucketRedis instance
	key := "my_bucket"
	rate := 10.0 // 10 requests per second
	lb := NewLeakyBucket(client, key, rate)

	// Test Allow method
	ctx := context.Background()

	// Test allowing requests
	toWait := lb.Allow(ctx)
	if toWait > 0 {
		time.Sleep(toWait)
		t.Errorf("Expected request to be allowed, but it was not")
	}

	toWait = lb.Allow(ctx)
	if toWait == 0 {
		t.Errorf("Expected request to be rate limited, but it was allowed")
	}

}
