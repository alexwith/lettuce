package storage

import (
	"strconv"
	"sync"
	"time"
)

var storage = make(map[string][]byte)
var timeouts = make(map[string]int64)
var storageMutex = &sync.RWMutex{}  // read-write lock
var timeoutsMutex = &sync.RWMutex{} // read-write lock

func Set(key string, value []byte, clearTimeout bool) {
	storageMutex.Lock()
	storage[key] = value
	storageMutex.Unlock()

	if clearTimeout {
		timeoutsMutex.Lock()
		delete(timeouts, key)
		timeoutsMutex.Unlock()
	}
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

func GetTimeout(key string) (int64, bool) {
	timeoutsMutex.RLock()
	value, present := timeouts[key]
	timeoutsMutex.RUnlock()

	return value, present
}

func Expire(key string, timeout int64) {
	timeoutsMutex.Lock()
	timeouts[key] = timeout
	timeoutsMutex.Unlock()
}

func ExpireIn(key string, milliseconds int64) {
	Expire(key, time.Now().UnixMilli()+milliseconds)
}

func Persist(key string) bool {
	_, keyPresent := Get(key)
	if !keyPresent {
		return false
	}

	_, timeoutPresent := GetTimeout(key)
	if !timeoutPresent {
		return false
	}

	timeoutsMutex.Lock()
	delete(timeouts, key)
	timeoutsMutex.Unlock()

	return true
}

func Increment(key string) (int, error) {
	value := storage[key]
	integer, err := strconv.Atoi(string(value))
	if err != nil {
		return -1, err
	}

	newInteger := integer + 1
	Set(key, []byte(strconv.Itoa(newInteger)), false)

	return newInteger, nil
}

func RegisterExpireTask() {
	ticker := time.NewTicker(250 * time.Millisecond)
	for {
		<-ticker.C

		currentTime := time.Now().UnixMilli()

		for key, timeout := range timeouts {
			if currentTime < timeout {
				continue
			}

			timeoutsMutex.Lock()
			delete(timeouts, key)
			timeoutsMutex.Unlock()

			Delete(key)
		}
	}
}
