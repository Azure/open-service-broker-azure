package generate

import "github.com/Azure/azure-service-broker/pkg/rand"

const (
	lowerAlphaChars = "abcdefghijklmnopqrstuvwxyz"
	upperAlphaChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars     = "0123456789"
)

var seededRand = rand.NewSeeded()
