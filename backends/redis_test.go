package backends

import (
	"github.com/garyburd/redigo/redis"
	"github.com/goibibo/norse/config"
	"github.com/stretchr/testify/assert"
	"testing"
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
	redisClient, _ := GetRedisClient(incrFun, decrFun, "redis")

	_, err := redisClient.Set("redisConfig", "key", "value")
	assert.NoError(t, err, "redis set error")

	value, errg := redisClient.Get("redisConfig", "key")
	assert.NoError(t, errg, "redis get error")
	assert.Equal(t, "value", value, "redis get value")

	delCount, errd := redisClient.Delete("redisConfig", "key")
	assert.NoError(t, errd, "redis delete error")
	assert.Equal(t, 1, delCount, "redis del value")
}

func Test_MGetSet(t *testing.T){
	redisClient, _ := GetRedisClient(incrFun, decrFun, "redis")

  val := map[string]interface{}{ "key1": "value1", "key2": "value2","key3": "value3"}
	_, err := redisClient.MSet("redisConfig", val)
	assert.NoError(t, err, "redis MSet error")

	value, errg := redisClient.MGet("redisConfig", "key1", "key2", "key3")
	assert.NoError(t, errg, "redis get error")
	assert.Equal(t, "value1", value[0], "redis get value")
  assert.Equal(t, "value2", value[1], "redis get value")
	assert.Equal(t, "value3", value[2], "redis get value")
}

func Test_Execute(t *testing.T) {
	redisClient, _ := GetRedisClient(incrFun, decrFun, "redis")

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
}
