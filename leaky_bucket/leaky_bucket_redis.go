package leaky_bucket_redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type LeakyBucketRedis struct {
	client *redis.Client
	key    string
	rate   float64 // Requests per second
}

func NewLeakyBucket(client *redis.Client, key string, rate float64) *LeakyBucketRedis {
	lb := &LeakyBucketRedis{
		client: client,
		key:    key,
		rate:   rate,
	}
	return lb
}

func (lb *LeakyBucketRedis) Allow(ctx context.Context) time.Duration {
	now := time.Now()

	// Leaky bucket algorithm
	script := `
		local ts  = tonumber(ARGV[1]) -- current time in seconds
		local cps = tonumber(ARGV[2]) -- capacity per second/rate
		local key = KEYS[1]

		-- remove tokens < min (older than now() -1s), keep only fresh tokens
		local min = ts - 1
		redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

		-- get the last token
		local last = redis.call('ZRANGE', key, -1, -1)
		local next = ts

		-- if there is a last token, calculate the next one based on the rate
		if type(last) == 'table' and #last > 0 then
			for key, value in pairs(last) do
				next = tonumber(value) + 1 / cps
				break -- break at first item
			end
		end

		if ts > next then
			-- the current ts is > than last+1/cps
			-- next = ts which is now
			next = ts
		end

		-- add the next token
		redis.call('ZADD', key, next, next)

		-- calculate the wait time based on difference between next and current ts
		local wait = next - ts
		return tostring(wait)

    `
	key := lb.key
	rate := lb.rate
	result, err := lb.client.Eval(ctx, script, []string{key}, now.Unix(), rate).Result()
	if err != nil {
		return 0
	}

	// Convert the result to duration in seconds
	wait, err := strconv.ParseFloat(result.(string), 64)
	if err != nil {
		return 0
	}

	// If wait time is negative, it means it's allowed, no need to wait
	if wait <= 0 {
		return 0
	}
	return time.Duration(wait * float64(time.Second))
}
