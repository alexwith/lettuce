package storage

import (
	"strconv"
	"sync"
	"time"
)

var storage = make(map[string][]byte)
var timeouts = make(map[string]int64)
var storageMutex = &sync.RWMutex{} // read-write lock
var timoutsMutex = &sync.RWMutex{} // read-write lock

func Set(key string, value []byte) {
	storageMutex.Lock()
	storage[key] = value
	storageMutex.Unlock()

	timoutsMutex.Lock()
	delete(timeouts, key)
	timoutsMutex.Unlock()
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

	timoutsMutex.RLock()
	currentExpireTime, present := timeouts[key]
	timoutsMutex.RUnlock()

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

	timeouts[key] = expireTime

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

func RegisterExpireTask() {
	ticker := time.NewTicker(250 * time.Millisecond)
	for {
		<-ticker.C

		currentTime := time.Now().UnixMilli()

		for key, expireTime := range timeouts {
			if currentTime < expireTime {
				continue
			}

			timoutsMutex.Lock()
			delete(timeouts, key)
			timoutsMutex.Unlock()

			Delete(key)
		}
	}
}
