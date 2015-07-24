package norse

import (
	"os"
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
var milliSecTimeout int

// Increment decrement functions
func incrFun(iKey string, incrementBy int64)error{return nil}
func decrFun(iKey string, decrementBy int64)error{return nil}

//type closureFunc func() (pool.Resource ,error)

var poolMap map[string]*pool.ResourcePool

// Redis connection struct
type RedisConn struct {
	redis.Conn
}

// All operations on redis from client happen through this struct
type RedisStruct struct{
	fIncr	func(string, int64)error
	fDecr	func(string, int64)error
	identifierkey	string
}

// Close redis conn
func (rConn *RedisConn) Close(){
	_ = rConn.Conn.Close()
}


func factory(key string, config map[string]string) (pool.Resource, error) {
	host := config["host"]
	port := config["port"]
	redisString := fmt.Sprintf("%s:%s", host, port)
	cli, err := redis.Dial("tcp", redisString)
	// select default db if not specified
	db, ok := config["db"]
	if ok{
		cli.Do("SELECT", db)
	}else{
		cli.Do("SELECT", 0)
	}
	if err != nil{
		// Write exit
		fmt.Println("Error in Redis Dial")
	}
	res := &RedisConn{cli}
	return res, nil
}


// Specify a factory function to create a connection,
// context and a timeout for connection to be created
func init() {
	// For each type in redis create corresponding pool
	poolMap = make(map[string]*pool.ResourcePool)
	milliSecTimeout = 5000
	redisConfigs, err := norse.LoadRedisConfig()
	if err != nil{
		os.Exit(1)
	}
	for key, config := range redisConfigs{
		factoryFunc := func(key string, config map[string]string) (pool.Factory) {
			return func()(pool.Resource,error){
				return factory(key, config)
			}
		}
		t := time.Duration(5000*time.Millisecond)
		poolMap[key] = pool.NewResourcePool(factoryFunc(key, config), 10, 100, t)
	}
}

// Your instance type for redis
func GetRedisClient(incr, decr func(string, int64)error, identifierKey string) (*RedisStruct) {
	return &RedisStruct{incr, decr, identifierKey}
}

// Execute, get connection from a pool 
// fetch and return connection to a pool.
func (r *RedisStruct) Execute(redisInstance string, cmd string, args ...interface{}) (interface{}, error) {
	ctx := context.TODO()
	// Get and set in our pool; for redis we use our own pool
	pool, ok := poolMap[redisInstance]
	// Increment and decrement counters using user specified functions.
	r.fIncr(r.identifierkey, 1)
	defer r.fDecr(r.identifierkey, 1)
	if ok{
		conn, _ := pool.Get(ctx)
		defer pool.Put(conn)
		return conn.(*RedisConn).Do(cmd, args...)
	}else{
		return nil, errors.New("Redis: instance Not found")
	}
}

// Redis Get,
func (r *RedisStruct) Get(redisInstance string, key string) (string, error){
	value, err := redis.String(r.Execute("flight", "GET", "a"))
	if err != nil{
		return value, nil
	}else{
		return "", err
	}
}

