package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Generic map
type M map[string]interface{}

var (
	Letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Digits  = []byte("0123456789")
	Symbols = []byte("!@#$%^&*-=_.?")
)
