package storage

import (
	"strconv"
	"sync"
	"time"
)

var storage = make(map[string][]byte)
var expiries = make(map[string]int64)
var storageMutex = &sync.RWMutex{} // read-write lock
var expireMutex = &sync.RWMutex{}  // read-write lock

func Set(key string, value []byte) {
	storageMutex.Lock()
	storage[key] = value
	storageMutex.Unlock()

	expireMutex.Lock()
	delete(expiries, key)
	expireMutex.Unlock()
}

func Get(key string) ([]byte, bool) {
	storageMutex.RLock()
	value, present := storage[key]
	storageMutex.RUnlock()

	return value, present
}

func Delete(key string) bool {
	_, present := Get(key)
	if !present {
		return false
	}

	storageMutex.Lock()
	delete(storage, key)
	storageMutex.Unlock()

	return true
}

func Expire(key string, seconds int, nx bool, xx bool, gt bool, lt bool) bool {
	expireTime := time.Now().UnixMilli() + int64(seconds*1000)

	_, keyPresent := Get(key)
	if !keyPresent {
		return false
	}

	expireMutex.RLock()
	currentExpireTime, present := expiries[key]
	expireMutex.RUnlock()
	if nx && present {
		return false
	}

	if xx && !present {
		return false
	}

	if gt && expireTime <= currentExpireTime {
		return false
	}

	if lt && expireTime >= currentExpireTime {
		return false
	}

	expiries[key] = expireTime

	return true
}

func Increment(key string) (int, error) {
	value := storage[key]
	integer, err := strconv.Atoi(string(value))
	if err != nil {
		return -1, err
	}

	newInteger := integer + 1
	Set(key, []byte(strconv.Itoa(newInteger)))

	return newInteger, nil
}
