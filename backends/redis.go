package backends

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

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
	fIncr func(string) error
	fDecr func(string) error
}

// Close redis conn
func (rConn *RedisConn) Close() {
	_ = rConn.Conn.Close()
}

// Callback function factory to vitess pool
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
		t := time.Duration(time.Duration(milliSecTimeout) * time.Millisecond)
		redisPoolMap[key] = pool.NewResourcePool(factoryFunc(key, config), 10, 100, t)
	}
}

// Your instance type for redis
func GetRedisClient(incr, decr func(string) error) (*RedisStruct, error) {
	if len(redisPoolMap) == 0 {
		return nil, errors.New("Redis Not Configured, please call Configure()")
	}
	return &RedisStruct{incr, decr}, nil
}

// Execute, get connection from a pool
// fetch and return connection to a pool.
func (r *RedisStruct) Execute(redisInstance string, cmd string, args ...interface{}) (interface{}, error) {
	// Get and set in our pool; for redis we use our own pool
	pool, ok := redisPoolMap[redisInstance]
	// Increment and decrement counters using user specified functions.
	if ok {
		conn, err := pool.Get(redisCtx)
		if err != nil {
			return nil, err
		}
		r.fIncr(redisInstance)
		defer r.fDecr(redisInstance)
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

// Redis HMGet,
func (r *RedisStruct) HMGet(redisInstance string, keys ...interface{}) ([]string, error) {
	values, err := redis.Strings(r.Execute(redisInstance, "HMGET", keys...))
	if err != nil {
		return []string{}, err
	} else {
		return values, err
	}
}

// Redis HMSet,
func (r *RedisStruct) HMSet(redisInstance string, key string, keyVapPair map[string]string) (bool, error) {
	_, err := r.Execute(redisInstance, "HMSET", redis.Args{key}.AddFlat(keyVapPair)...)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// Redis SetEx
func (r *RedisStruct) Setex(redisInstance string, key string, duration int, value interface{}) (string, error) {
	_, err := r.Execute(redisInstance, "SETEX", key, duration, value)
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

// Redis MSet
func (r *RedisStruct) MSet(redisInstance string, keyVapPair map[string]interface{}) (bool, error) {
	_, err := r.Execute(redisInstance, "MSET", redis.Args{}.AddFlat(keyVapPair)...)
	if err != nil {
		return false, err
	}
	return true, nil
}

// redis SMEMBERS
func (r *RedisStruct) Smembers(redisInstance string, key string) ([]string, error) {
	val, err := redis.Values(r.Execute(redisInstance, "SMEMBERS", key))
	if err != nil {
		return []string{}, err
	}
	s := make([]string, len(val))
	//Convert array of Bytes to array of string
	for i, item := range val {
		s[i] = string(item.([]byte))
	}
	return s, nil
}

// redis SADD
func (r *RedisStruct) SAdd(redisInstance string, key string, values ...interface{}) (bool, error) {
	_, err := r.Execute(redisInstance, "SADD", redis.Args{}.Add(key).AddFlat(values)...)
	if err != nil {
		return false, err
	}
	return true, nil
}

// redis SREM
func (r *RedisStruct) SRem(redisInstance string, key string, value string) (bool, error) {
	_, err := r.Execute(redisInstance, "SREM", key, value)
	if err != nil {
		return false, err
	}
	return true, nil
}

// redis SISMEMBER
func (r *RedisStruct) Sismember(redisInstance string, key string, member string) (bool, error) {
	val, err := r.Execute(redisInstance, "SISMEMBER", key, member)
	if err != nil {
		return false, err
	}
	// val is interface; trying to convert to int64
	return val.(int64) != 0, nil
}

func (r *RedisStruct) Sismembers(redisInstance string, key string, members []string) ([]bool, error) {

	// Get and set in our pool; for redis we use our own pool
	pool, ok := redisPoolMap[redisInstance]
	// Increment and decrement counters using user specified functions.
	if ok {
		conn, err := pool.Get(redisCtx)
		if err != nil {
			return nil, err
		}
		r.fIncr(redisInstance)
		defer r.fDecr(redisInstance)
		defer pool.Put(conn)

		for _, member := range members {
			conn.(*RedisConn).Send("SISMEMBER", key, member)
		}
		conn.(*RedisConn).Flush()

		results := make([]bool, 0, len(members))
		for _, _ = range members {
			res, err := conn.(*RedisConn).Receive()
			if err != nil {
				return nil, err
			}
			val := res.(int64) != 0
			results = append(results, val)
		}
		return results, nil

	} else {
		return nil, errors.New("Redis: instance Not found")
	}
}

func (r *RedisStruct) Delete(redisInstance string, keys ...interface{}) (int, error) {
	value, err := redis.Int(r.Execute(redisInstance, "DEL", keys...))
	if err != nil {
		return -1, err
	}
	return value, nil
}

//Redis LPUSH
func (r *RedisStruct) LPush(redisInstance string, key string, values ...interface{}) (bool, error) {
	inter := []interface{}{key}
	inter = append(inter, values...)
	_, err := r.Execute(redisInstance, "LPUSH", inter...)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Redis LRANGE
func (r *RedisStruct) LRange(redisInstance string, key string, start, end int) ([]string, error) {
	list, err := redis.Strings(r.Execute(redisInstance, "LRANGE", key, start, end))
	if err != nil {
		return []string{}, err
	}
	return list, nil
}

//Redis ZRANGE
func (r *RedisStruct) ZRange(redisInstance string, key string, start, end int) ([]string, error) {
	value, err := redis.Strings(r.Execute(redisInstance, "ZRANGE", key, start, end))
	if err != nil {
		return []string{}, err
	}
	return value, nil
}

//Redis ZRANGE with scores
func (r *RedisStruct) ZRangeWithScores(redisInstance string, key string, start, end int) ([]string, error) {
	value, err := redis.Strings(r.Execute(redisInstance, "ZRANGE", key, start, end, "withscores"))
	if err != nil {
		return []string{}, err
	}
	return value, nil
}

//Redis ZADD
func (r *RedisStruct) ZAdd(redisInstance string, key string, score int, member string) (bool, error) {
	_, err := r.Execute(redisInstance, "ZADD", key, score, member)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Redis ZREM
func (r *RedisStruct) ZRem(redisInstance string, key string, member string) (bool, error) {
	_, err := r.Execute(redisInstance, "ZREM", key, member)
	if err != nil {
		return false, err
	}
	return true, nil
}
