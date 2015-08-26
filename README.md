# norse

See config file in gist [here](https://gist.github.com/pranjal5215/cb8977317023f0f2ba40)

####Code:

	package main

	import (
		"fmt"
		"github.com/goibibo/norse/backends"
		"github.com/goibibo/norse/config"
	)

	// Increment decrement functions
	func incrFun(iKey string)error{return nil}
	func decrFun(iKey string)error{return nil}

	func main(){
		// See gist for sample config
		// Save config.json and set path to that path
		config.Configure(path)
		backends.Configure()

		// How to use redis,
		redisClient, _ := norse.GetRedisClient(incrFun, decrFun)
		value, _ := redisClient.Get("redisConfig", "key"))
		fmt.Println(value)


		// How to use memcache,
		memcacheClient, _ := norse.GetMemcacheClient(incrFun, decrFun)
		value, _ := memcacheClient.Get("configType", "key"))
		fmt.Println(value)


		// How to use MySQL,
		mysqlClient :=norse.GetMysqlClient(incrFunc,decrFunc,"mysql")
		// "mysql" database name;
		mysqlmap := mysqlClient.Select("select * from mytable")	
		fmt.Println(mysqlmap)

	}
