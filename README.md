# norse

####Code:

	package main

	import (
		"fmt"
		"github.com/goibibo/norse"
	)

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
