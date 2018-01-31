package cachetool

import (
	"EventFlow/common/tool/parametertool"
	"log"
	"sync"
	"time"
)

type Cache struct {
	CacheKey      string `yaml:cachekey`
	TimeoutSecond int    `yaml:timeoutsecond`
}

type cacheData struct {
	timer *time.Timer
	data  interface{}
}

var mutex sync.Mutex
var cacheMap map[string]cacheData

func init() {
	mutex = sync.Mutex{}
	cacheMap = make(map[string]cacheData)
}

func GetCacheByParameter(cache Cache, parameters *map[string]interface{}) (cacheValue interface{}, existed bool) {
	key := parametertool.ReplaceWithParameter(&cache.CacheKey, parameters)

	if key == "" {
		return nil, false
	}

	return GetCache(key)
}

func GetCache(cacheKey string) (cacheValue interface{}, existed bool) {

	if cacheKey == "" {
		return nil, false
	}

	mutex.Lock()

	cacheData, existed := cacheMap[cacheKey]

	if existed {
		cacheValue = cacheData.data
	} else {
		log.Printf("[cachetool] cache data doesn,t exist", cacheKey)
	}

	mutex.Unlock()

	return
}

func CreateCacheByParameter(cache Cache, cacheValue interface{}, parameters *map[string]interface{}) {
	key := parametertool.ReplaceWithParameter(&cache.CacheKey, parameters)

	if key != "" {
		CreateCache(key, cache.TimeoutSecond, cacheValue)
	}
}

func CreateCache(cacheKey string, timeoutSeconds int, cacheValue interface{}) {

	duration := time.Second * time.Duration(timeoutSeconds)

	mutex.Lock()

	existedCacheData, existed := cacheMap[cacheKey]

	if existed {
		existedCacheData.timer.Reset(duration)
		existedCacheData.data = cacheValue
	} else {
		timer := time.AfterFunc(duration, func() {
			delete(cacheMap, cacheKey)
		})
		cacheMap[cacheKey] = cacheData{data: cacheValue, timer: timer}
	}

	mutex.Unlock()
}
