package main

import "errors"

var errKeyNotFount = errors.New("key not found")

type api interface {
	get(key string) (string, error)
	put(key string, value string) error
	del(key string) error
}
