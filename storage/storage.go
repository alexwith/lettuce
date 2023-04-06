package storage

import (
	"strconv"
)

var storage map[string][]byte = make(map[string][]byte)

func Set(key string, value []byte) {
	storage[key] = value
}

func Get(key string) []byte {
	return storage[key]
}

func Increment(key string) int {
	value := storage[key]
	integer, _ := strconv.Atoi(string(value))

	newInteger := integer + 1
	Set(key, []byte(strconv.Itoa(newInteger)))

	return newInteger
}
