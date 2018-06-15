package crypto

import (
	"errors"
	"sync"
)

// Codec is an interface to be implemented by any type that can encrypt and
// decrypt values
type Codec interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

var globalCodec Codec
var globalCodecMutex sync.RWMutex

// InitializeGlobalCodec may be called once and only once to inject an object
// that implements the Codec interface into a "global" (package-scoped)
// variable. Additional invocations of this function will return an error.
func InitializeGlobalCodec(codec Codec) error {
	globalCodecMutex.Lock()
	defer globalCodecMutex.Unlock()
	if globalCodec != nil {
		return errors.New("Global codec may not be re-initialized")
	}
	globalCodec = codec
	return nil
}

// Encrypt encrypts the provided bytes using the globally configured codec
func Encrypt(bytes []byte) ([]byte, error) {
	globalCodecMutex.RLock()
	defer globalCodecMutex.RUnlock()
	if globalCodec == nil {
		return nil, errors.New("No global codec has been configured")
	}
	return globalCodec.Encrypt(bytes)
}

// Decrypt decrypts the provided bytes using the globally configured codec
func Decrypt(bytes []byte) ([]byte, error) {
	globalCodecMutex.RLock()
	defer globalCodecMutex.RUnlock()
	if globalCodec == nil {
		return nil, errors.New("No global codec has been configured")
	}
	return globalCodec.Decrypt(bytes)
}
