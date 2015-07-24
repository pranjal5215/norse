package main

import (
	"fmt"
	"time"
	"errors"
	"golang.org/x/net/context"

	"github.com/goibibo/norse"
	"github.com/goibibo/hammerpool/pool"
	"github.com/garyburd/redigo/redis"
)


// Get redis config 
var redisConfigs map[string]map[string]string

// redis timeout
var milliSecTimeout = 5000

// Increment decrement functions
var incrFun func
var decrFun func

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
	redisConfigs, err := norse.LoadRedisConfig()
	for key, config := range redisConfigs{
		factory := func() (hammerpool.Resource, error) {
			host := config["host"]
			port := config["port"]
			redisString := fmt.Sprintf("%s:%s", host, port)
			cli, err := redis.Dial("tcp", redisString)
			// select default db if not specified
			db, ok := config["db"]
			if ok{
				cli.Do("SELECT", config)
			}else{
				cli.Do("SELECT", 0)
			}
			if err != nil{
				// Write exit
				fmt.Println("Error in Redis Dial")
			}
			res := &RedisStruct{cli}
			return res, nil
		}
		t := time.Duration(milliSecTimeout*time.Millisecond)
		poolMap[key] = hammerpool.NewResourcePool(factory, 10, t)
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
	defer pool.Put(conn)
	// Increment and decrement counters using user specified functions.
	incrFun()
	defer dncrFun()
	if ok{
		conn, err := pool.Get(ctx)
	}else{
		return nil, errors.New("Redis: instance Not found")
	}
	if err != nil {
		return nil, err
	}
	return conn.(*RedisStruct).Do(cmd, args...)
}

// Set incr decr functions
func (r *RedisStruct) SetIncrDecrFunctions(incr func, decr func) (string, error){
	incrFun = incr
	decrFun = decr
}

// Redis Get,
func (r *RedisStruct) Get(redisInstance string, key string) (string, error){
	value, err := redis.String(rStruct.Execute("flight", "GET", "a", Incr, Decr))
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
