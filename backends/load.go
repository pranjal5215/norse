package backends

import (
	"sync"
)

var doOnce *sync.Once

func init() {
	doOnce = &sync.Once{}
}

func Configure() {
	doOnce.Do(doConfigure)
}

func doConfigure() {
	configureRedis()
	configureMemcache()
	configureMySql()
}
