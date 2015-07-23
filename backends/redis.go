package main

import (
	"fmt"
	"time"
	"errors"
	"golang.org/x/net/context"
	"github.com/goibibo/hammerpool/pool"
	"github.com/garyburd/redigo/redis"
)



var poolMap map[string]*pool.ResourcePool

type RedisStruct struct {
	redis.Conn
}

// Close redis conn
func (rConn *RedisStruct) Close(){
	_ = rConn.Conn.Close()
}

// Specify a factory function to create a connection,
// context and a timeout for connection to be created
func init() {
	// For each type in redis create corresponding pool
	for key, config := range redisConfigs{
		factory := func() (pool.Resource, error) {
			cli, err := redis.Dial("tcp", "127.0.0.1:6379")
			// select db if specified
			db, ok := config["db"]
			if ok{
				cli.Do("SELECT", config)
			}else{
				cli.Do("SELECT", 0)
			}
			if err != nil{
				fmt.Println("Error in Redis Dial")
			}
			res := &RedisStruct{cli}
			return res, nil
		}
		t := time.Duration(5000*time.Millisecond)
		poolMap[key] = pool.NewResourcePool(factory, 10, 100, t)
	}
}

// Your instance type for redis
func GetMyRedis() (*RedisStruct) {
	return &RedisStruct{}
}

// Execute, get connection from a pool 
// fetch and return connection to a pool.
func (r *RedisStruct) Execute(redisInstance string, cmd string, args ...interface{}) (interface{}, error) {
	ctx := context.TODO()
	pool, ok := poolMap[redisInstance]
	if ok{
		conn, err := pool.Get(ctx)
	}else{
		return nil, errors.New("Redis: instance Not found")
	}
	if err != nil {
		return nil, err
	}
	defer pool.Put(conn)
	return conn.(*RedisStruct).Do(cmd, args...)
}

// Redis Get,
func Get(redisInstance string, key string) (string, error){
	value, err := redis.String(rStruct.Execute("flight", "GET", "a"))
	if err != nil{
		return value, nil
	}else{
		return nil, err
	}
}

// How to use redis,
func main(){
	rStruct := GetMyRedis()
	value, _ := redis.String(rStruct.Execute("flight", "GET", "a"))
	fmt.Println(value)
	value, _ = redis.String(rStruct.Execute("flight", "GET", "key"))
	fmt.Println(value)
}
