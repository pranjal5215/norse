package backends

import (
	"errors"
	"golang.org/x/net/context"
	"os"
	"time"

	"github.com/goibibo/gomemcache/memcache"
	"github.com/goibibo/hammerpool/pool"
	"github.com/goibibo/norse/config"
)

var (
	// Get memcache config
	memcacheConfigs map[string]map[string]string

	// memcache timeout
	memMilliSecTimeout int

	// type closureFunc func() (pool.Resource ,error)
	memPoolMap map[string]*pool.ResourcePool

	// context var to vitess pool
	memCtx context.Context
)

// Memcache connection struct
type MemcacheConn struct {
	*memcache.Client
}

// All operations on memcache from client happen through this struct
type MemcacheStruct struct {
	fIncr         func(string, int64) error
	fDecr         func(string, int64) error
	identifierkey string
}

// Close memcache conn
func (mConn *MemcacheConn) Close() {
}

// Callback function factory to vitess pool`
func memcacheFactory(key string, server []string) (pool.Resource, error) {
	conn := memcache.New(server...)
	memConn := &MemcacheConn{conn}
	return memConn, nil
}

// Specify a factory function to create a connection,
// context and a timeout for connection to be created
func init() {
	// For each type in memcache create corresponding pool
	memCtx = context.Background()
	memPoolMap = make(map[string]*pool.ResourcePool)
	milliSecTimeout = 5000
	memcacheConfigs, err := config.LoadMemcacheConfig()
	if err != nil {
		os.Exit(1)
	}
	for key, config := range memcacheConfigs {
		factoryFunc := func(key string, config []string) pool.Factory {
			return func() (pool.Resource, error) {
				return memcacheFactory(key, config)
			}
		}
		t := time.Duration(5000 * time.Millisecond)
		memPoolMap[key] = pool.NewResourcePool(factoryFunc(key, config), 10, 100, t)
	}
}

// Your instance type for memcache
func GetMemcacheClient(incr, decr func(string, int64) error, identifierKey string) *MemcacheStruct {
	return &MemcacheStruct{incr, decr, identifierKey}
}

// Memcache Get,
func (m *MemcacheStruct) Get(memcacheInstance string, key string) (string, error) {
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := memPoolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok {
		conn, _ := pool.Get(memCtx)
		defer pool.Put(conn)
		value, err := conn.(*MemcacheConn).Get(key)
		if err != nil {
			return "", err
		}
		return string(value.Value), err
	} else {
		return "", errors.New("Memcache: instance Not found")
	}
}

// Memcache Set,
func (m *MemcacheStruct) Set(memcacheInstance string, key string, value string) (string, error) {
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := memPoolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok {
		conn, _ := pool.Get(memCtx)
		defer pool.Put(conn)
		byteArr := []byte(value)
		conn.(*MemcacheConn).Set(&memcache.Item{Key: key, Value: byteArr})
		return "", nil
	} else {
		return "", errors.New("Memcache: instance Not found")
	}
}

func (m *MemcacheStruct) Setex(memcacheInstance string, key string, duration int, val string) (bool, error) {
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := memPoolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok {
		conn, _ := pool.Get(memCtx)
		defer pool.Put(conn)
		resp := conn.(*MemcacheConn).Set(&memcache.Item{Key: key, Value: []byte(val)})
		if resp != nil {
			return false, resp
		} else {
			err := conn.(*MemcacheConn).Touch(key, int32(duration))
			if err != nil {
				return false, err
			}
			return true, nil
		}
	} else {
		return false, errors.New("Memcache: instance Not found")
	}
}

func (m *MemcacheStruct) Expire(memcacheInstance string, key string, duration int) (bool, error) {
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := memPoolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok {
		conn, _ := pool.Get(memCtx)
		defer pool.Put(conn)
		err := conn.(*MemcacheConn).Touch(key, int32(duration))
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return true, errors.New("Memcache: instance Not found")
	}
}
