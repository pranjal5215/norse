package backends

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/goibibo/hammerpool/pool"
	"github.com/goibibo/norse/config"
)

var (
	// Get redis config
	redisConfigs map[string]map[string]string

	// redis timeout
	milliSecTimeout int

	// type closureFunc func() (pool.Resource ,error)
	redisPoolMap map[string]*pool.ResourcePool

	// context var to vitess pool
	redisCtx context.Context
)

// Redis connection struct
type RedisConn struct {
	redis.Conn
}

// All operations on redis from client happen through this struct
type RedisStruct struct {
	fIncr         func(string, int64) error
	fDecr         func(string, int64) error
	identifierkey string
}

// Close redis conn
func (rConn *RedisConn) Close() {
	_ = rConn.Conn.Close()
}

// Callback function factory to vitess pool`
func redisFactory(key string, config map[string]string) (pool.Resource, error) {
	host := config["host"]
	port := config["port"]
	redisString := fmt.Sprintf("%s:%s", host, port)
	cli, err := redis.Dial("tcp", redisString)
	// select default db if not specified
	db, ok := config["db"]
	if ok {
		cli.Do("SELECT", db)
	} else {
		cli.Do("SELECT", 0)
	}
	if err != nil {
		// Write exit
		fmt.Println("Error in Redis Dial")
	}
	res := &RedisConn{cli}
	return res, nil
}

// Specify a factory function to create a connection,
// context and a timeout for connection to be created
func configureRedis() {
	// For each type in redis create corresponding pool
	redisCtx = context.Background()
	redisPoolMap = make(map[string]*pool.ResourcePool)
	milliSecTimeout = 5000
	redisConfigs, err := config.LoadRedisConfig()
	if err != nil {
		os.Exit(1)
	}
	for key, config := range redisConfigs {
		factoryFunc := func(key string, config map[string]string) pool.Factory {
			return func() (pool.Resource, error) {
				return redisFactory(key, config)
			}
		}
		t := time.Duration(5000 * time.Millisecond)
		redisPoolMap[key] = pool.NewResourcePool(factoryFunc(key, config), 10, 100, t)
	}
}

// Your instance type for redis
func GetRedisClient(incr, decr func(string, int64) error, identifierKey string) *RedisStruct, error {
	if (len(redisPoolMap) == 0){
		return nil, errors.New("Redis Not Configured, please call Configure()")
	}
	return &RedisStruct{incr, decr, identifierKey}, nil
}

// Execute, get connection from a pool
// fetch and return connection to a pool.
func (r *RedisStruct) Execute(redisInstance string, cmd string, args ...interface{}) (interface{}, error) {
	// Get and set in our pool; for redis we use our own pool
	pool, ok := redisPoolMap[redisInstance]
	// Increment and decrement counters using user specified functions.
	r.fIncr(r.identifierkey, 1)
	defer r.fDecr(r.identifierkey, 1)
	if ok {
		conn, _ := pool.Get(redisCtx)
		defer pool.Put(conn)
		return conn.(*RedisConn).Do(cmd, args...)
	} else {
		return nil, errors.New("Redis: instance Not found")
	}
}

// Redis Get,
func (r *RedisStruct) Get(redisInstance string, key string) (string, error) {
	value, err := redis.String(r.Execute(redisInstance, "GET", key))
	if err != nil {
		return "", nil
	} else {
		return value, err
	}
}

// Redis Set,
func (r *RedisStruct) Set(redisInstance string, key string, value interface{}) (string, error) {
	_, err := r.Execute(redisInstance, "SET", key, value)
	if err != nil {
		return "", err
	} else {
		return "", nil
	}
}

// Redis MGet
func (r *RedisStruct) MGet(redisInstance string, keys ...interface{}) ([]string, error) {
	values, err := redis.Strings(r.Execute(redisInstance, "MGET", keys...))
	if err != nil {
		return []string{}, err
	}
	return values, nil
}
