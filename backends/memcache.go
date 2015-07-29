package norse

import (
	"os"
	"fmt"
	"time"
	"errors"
	"golang.org/x/net/context"

	"github.com/goibibo/norse"
	"github.com/goibibo/hammerpool/pool"
	"github.com/goibibo/gomemcache/memcache"
)

var (
	// Get memcache config 
	memcacheConfigs map[string]map[string]string

	// memcache timeout
	milliSecTimeout int

	// type closureFunc func() (pool.Resource ,error)
	poolMap map[string]*pool.ResourcePool

	// context var to vitess pool
	ctx    context.Context
)

// Memcache connection struct
type MemcacheConn struct {
	*memcache.Client
}

// All operations on memcache from client happen through this struct
type MemcacheStruct struct{
	fIncr	func(string, int64)error
	fDecr	func(string, int64)error
	identifierkey	string
}

// Close memcache conn
func (rConn *MemcacheConn) Close(){
}

// Callback function factory to vitess pool`
func factory(key string, config []string) (pool.Resource, error) {
	res := memcache.New(server...)
	return res, nil
}

// Specify a factory function to create a connection,
// context and a timeout for connection to be created
func init() {
	// For each type in memcache create corresponding pool
	ctx = context.Background()
	poolMap = make(map[string]*pool.ResourcePool)
	milliSecTimeout = 5000
	memcacheConfigs, err := norse.LoadMemcacheConfig()
	if err != nil{
		os.Exit(1)
	}
	for key, config := range memcacheConfigs{
		factoryFunc := func(key string, config []string) (pool.Factory) {
			return func()(pool.Resource,error){
				return factory(key, config)
			}
		}
		t := time.Duration(5000*time.Millisecond)
		poolMap[key] = pool.NewResourcePool(factoryFunc(key, config), 10, 100, t)
	}
}

// Your instance type for memcache
func GetMemcacheClient(incr, decr func(string, int64)error, identifierKey string) (*MemcacheStruct) {
	return &MemcacheStruct{incr, decr, identifierKey}
}

// Memcache Get,
func (m *MemcacheStruct) Get(memcacheInstance string, key string) (string, error){
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := poolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok{
		conn = pool.Get(ctx)
		defer pool.Put(conn)
		value, err := conn.Get(key)
		if err != nil{
			return "", err
		}
		return value, err
	}else{
		return nil, errors.New("Memcache: instance Not found")
	}
}

// Memcache Set,
func (m *MemcacheStruct) Set(memcacheInstance string, key string, value string) (string, error){
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := poolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok{
		conn = pool.Get(ctx)
		defer pool.Put(conn)
		byteArr := []byte(value)
		conn.Set(&memcache.Item{Key:key, Value:byteArr})
		return true, nil
	}else{
		return nil, errors.New("Memcache: instance Not found")
	}
}

func (m *MemcacheStruct) Setex(memcacheInstance string, key string, duration int, val string) (bool, error) {
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := poolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok{
		resp := conn.Set(&memcache.Item{Key: key, Value: []byte(val)})
		if resp != nil {
			return false, erm
		} else {
			err := conn.Touch(key, int32(duration))
			if err != nil {
				return false, err
			}
		}
	}else{
		return nil, errors.New("Memcache: instance Not found")
	}
}

func (m *MemcacheStruct) Expire(memcacheInstance string, key string, duration int) (bool, error) {
	// Get and set in our pool; for memcache we use our own pool
	pool, ok := poolMap[memcacheInstance]
	// Increment and decrement counters using user specified functions.
	m.fIncr(m.identifierkey, 1)
	defer m.fDecr(m.identifierkey, 1)
	if ok{
		conn = pool.Get(ctx)
		defer pool.Put(conn)
		err := conn.Touch(key, int32(duration))
		if err!=nil{
			return "", err
		}
		return true, nil
	}else{
		return nil, errors.New("Memcache: instance Not found")
	}
}
