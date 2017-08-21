package generate

import (
	"math/rand"
	"sync"
	"time"
)

const (
	lowerAlphaChars = "abcdefghijklmnopqrstuvwxyz"
	upperAlphaChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars     = "0123456789"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var seededRandMutex = &sync.Mutex{}
