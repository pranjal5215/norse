package backends

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/goibibo/norse/config"
	"github.com/stretchr/testify/assert"
)

var redisCount int

func incrFun(iKey string) error { redisCount++; return nil }
func decrFun(iKey string) error { redisCount--; return nil }

func init() {

	redisCount = 0

	config.Configure("./test_config.json")
	Configure()
}

func Test_GetSet(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun)

	_, err := redisClient.Set("redisConfig", "key", "value")
	assert.NoError(t, err, "redis set error")

	value, errg := redisClient.Get("redisConfig", "key")
	assert.NoError(t, errg, "redis get error")
	assert.Equal(t, "value", value, "redis get value")

	delCount, errd := redisClient.Delete("redisConfig", "key")
	assert.NoError(t, errd, "redis delete error")
	assert.Equal(t, 1, delCount, "redis del value")

	redisClient.Execute("redisConfig", "FLUSHDB")
}

func Test_MGetSet(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun)

	val := map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}
	_, err := redisClient.MSet("redisConfig", val)
	assert.NoError(t, err, "redis MSet error")

	value, errg := redisClient.MGet("redisConfig", "key1", "key2", "key3")
	assert.NoError(t, errg, "redis get error")
	assert.Equal(t, "value1", value[0], "redis get value")
	assert.Equal(t, "value2", value[1], "redis get value")
	assert.Equal(t, "value3", value[2], "redis get value")

	redisClient.Execute("redisConfig", "FLUSHDB")
}

func Test_HMGetSet(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun)
	val := map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}
        _, err := redisClient.HMSet("redisConfig", "key", val)
	assert.NoError(t, err, "redis HMSet error")

	value, errg := redisClient.HMGet("redisConfig", "key", "key1")
	assert.NoError(t, errg, "redis get error")
	assert.Equal(t, "value1", value[0], "redis get value")

	redisClient.Execute("redisConfig", "FLUSHDB")
}

func Test_Execute(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun)

	_, err := redisClient.Execute("redisConfig", "SET", "key", "value")
	assert.NoError(t, err, "redis set error")

	value, errg := redisClient.Execute("redisConfig", "GET", "key")
	assert.NoError(t, errg, "redis get error")
	vals, errs := redis.String(value, nil)
	assert.NoError(t, errs, "redis string error")
	assert.Equal(t, "value", vals, value, "redis get value")

	delCount, errd := redisClient.Execute("redisConfig", "DEL", "key")
	assert.NoError(t, errd, "redis delete error")
	dcI, erri := redis.Int(delCount, nil)
	assert.NoError(t, erri, "redis int error")
	assert.Equal(t, 1, dcI, "redis del value")

	redisClient.Execute("redisConfig", "FLUSHDB")
}

func Test_Sets(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun)

	saddb, err := redisClient.SAdd("redisConfig", "set", "ele1", "ele2", "ele3", "ele4")
	assert.Equal(t, true, saddb, "sadd")
	assert.NoError(t, err, "sadd error")

	set, err1 := redisClient.Sismember("redisConfig", "set", "ele1")
	assert.Equal(t, true, set, "sismember ")
	assert.NoError(t, err1, "sadd error")

	set, err1 = redisClient.Sismember("redisConfig", "set", "ele2")
	assert.Equal(t, true, set, "sismember")
	assert.NoError(t, err1, "sadd error")

	set, err1 = redisClient.Sismember("redisConfig", "set", "ele3")
	assert.Equal(t, true, set, "sismember")
	assert.NoError(t, err1, "sadd error")

	set, err1 = redisClient.Sismember("redisConfig", "set", "ele4")
	assert.Equal(t, true, set, "sismember")
	assert.NoError(t, err1, "sadd error")

	sets, err2 := redisClient.Sismembers("redisConfig", "set", []string{"ele1", "ele2", "ele3", "ele4"})
	assert.Equal(t, []bool{true, true, true, true}, sets, "sismembers")
	assert.NoError(t, err2, "sadd error")

	redisClient.Execute("redisConfig", "FLUSHDB")
}

func Test_Sets2(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun)

	saddb, err := redisClient.SAdd("redisConfig", "set", "ele1", "ele2", "ele3", "ele4")
	assert.Equal(t, true, saddb, "sadd")
	assert.NoError(t, err, "sadd error")

	set, err1 := redisClient.SRem("redisConfig", "set", "ele1")
	assert.Equal(t, true, set, "srem ")
	assert.NoError(t, err1, "srem error")

	set, err1 = redisClient.Sismember("redisConfig", "set", "ele1")
	assert.Equal(t, false, set, "sismember ")
	assert.NoError(t, err1, "sadd error")

	sets, err2 := redisClient.Sismembers("redisConfig", "set", []string{"ele2", "ele3", "ele4"})
	assert.Equal(t, []bool{true, true, true}, sets, "sismembers")
	assert.NoError(t, err2, "sadd error")

	saddb2, err := redisClient.SAdd("redisConfig", "set", "ele1", "ele2", "ele3", "ele4")
	assert.Equal(t, true, saddb2, "sadd")
	assert.NoError(t, err, "sadd error")

	set2, err1 := redisClient.Sismember("redisConfig", "set", "ele1")
	assert.Equal(t, true, set2, "sismember ")
	assert.NoError(t, err1, "sadd error")

	redisClient.Execute("redisConfig", "FLUSHDB")
}
