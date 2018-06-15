package generate

import "open-service-broker-azure/pkg/rand"

const (
	lowerAlphaChars = "abcdefghijklmnopqrstuvwxyz"
	upperAlphaChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars     = "0123456789"
)

var seededRand = rand.NewSeeded()
