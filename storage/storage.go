package storage

import (
	"strconv"
	"sync"
)

var storage = make(map[string][]byte)
var storageMutex = &sync.RWMutex{} // read-write lock

func Set(key string, value []byte) {
	storageMutex.Lock()
	storage[key] = value
	storageMutex.Unlock()
}

func Get(key string) ([]byte, bool) {
	storageMutex.RLock()

	value, present := storage[key]

	storageMutex.RUnlock()

	return value, present
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
