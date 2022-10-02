package util

import "github.com/coocood/freecache"

// 初始化Cache
var cache = freecache.NewCache(100 * 1024 * 1024)

// CacheModel 定一个 model struct
type CacheModel struct {
	Key           string
	Value         string
	ExpireSeconds int //过期时间-s
}

func GetCache(key string) (string, bool) {
	value, err := cache.Get([]byte(key))
	if err != nil {
		return "", false
	}
	return string(value), true
}

// GetOrSetCache GetOrSet
// 如果没有就存入新的key
func GetOrSetCache(free CacheModel) (string, bool) {
	retValue, err := cache.GetOrSet([]byte(free.Key), []byte(free.Value), free.ExpireSeconds)
	if err != nil {
		return "", false
	}
	return string(retValue), true
}

// SetCache Set
func SetCache(free CacheModel) bool {
	err := cache.Set([]byte(free.Key), []byte(free.Value), free.ExpireSeconds)
	if err != nil {
		return false
	}
	return true
}

// SetAndGetCache SetAndGet
func SetAndGetCache(free CacheModel) (string, bool) {
	retValue, found, err := cache.SetAndGet([]byte(free.Key), []byte(free.Value), free.ExpireSeconds)
	if err != nil {
		return "", false
	}
	return string(retValue), found
}

// TtlCache 更新key的过期时间--如果这个key不存在会返回错误
func TtlCache(free CacheModel) bool {
	err := cache.Touch([]byte(free.Key), free.ExpireSeconds)
	if err != nil {
		return false
	}
	return true
}

// DelCache 删除---
func DelCache(free CacheModel) bool {
	return cache.Del([]byte(free.Key))
}
