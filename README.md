# Leaky Bucket Redis

## Overview
The `leaky_bucket_redis.go` is a leaky bucket implementation in Go that utilizes Redis's Lua script. It is designed to work well with distributed systems, providing atomic operations for managing rate limiting and throttling.

## Installation
To use the `leaky_bucket_redis.go` in your Go project, you need to have Redis installed and running. You can install Redis by following the instructions provided in the official Redis documentation.

Once Redis is installed, you can add the `leaky_bucket_redis.go` file to your project's source code.

## Usage
To use the leaky bucket implementation, follow these steps:

1. Import the necessary packages in your Go file:
   ```go
   import (
       "github.com/redis/go-redis/v9"
       "github.com/alibazlamit/leaky_bucket_redis"
   )

### LeakyBucketRedis Allow Method

The `Allow` method in the `LeakyBucketRedis` struct implements the Leaky Bucket algorithm. It takes a `context.Context` as input and returns A `time.Duration` representing the wait time if the request is not allowed.

The method works as follows:

- It obtains the current time and defines a Lua script for the Leaky Bucket algorithm.
- It assigns the Redis key and the rate of token consumption from the bucket.
- It executes the Lua script on the Redis server using the `Eval` method of the Redis client.
- It parses the result of the Lua script to a float64 representing the wait time in seconds.
- If the wait time is less than or equal to 0, the method returns wait time of 0, indicating that the request is allowed.
- If the wait time is greater than 0, the method returns the calculated wait time, indicating that the request is not allowed.

This implementation allows you to control the rate of requests in your application using the Leaky Bucket algorithm.


If you have any suggestions or find this implementation useful, feel free to drop me a star!
