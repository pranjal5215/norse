# norse

####Code:

	package main

	import (
		"fmt"
		"github.com/goibibo/norse"
	)

	// Increment decrement functions
	func incrFun(iKey string, incrementBy int64)error{return nil}
	func decrFun(iKey string, decrementBy int64)error{return nil}

	// How to use redis,
	func main(){
		redisClient := norse.GetRedisClient(incrFun, decrFun, "redis")
		value, _ := redis.String(redisClient.Execute("flight", "GET", "a"))
		fmt.Println(value)
		value, _ = redis.String(redisClient.Execute("flight", "GET", "key"))
		fmt.Println(value)
		value, _ = redis.String(redisClient.Execute("bus", "GET", "hhh"))
		fmt.Println(value)
		value, _ = redis.String(redisClient.Execute("bus", "GET", "xxx"))
		fmt.Println(value)
	}
