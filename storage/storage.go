package storage

import (
	"strconv"
)

var storage map[string][]byte = make(map[string][]byte)

func Set(key string, value []byte) {
	storage[key] = value
}

func Get(key string) ([]byte, bool) {
	value, present := storage[key]
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
