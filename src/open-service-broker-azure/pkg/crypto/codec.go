package crypto

// Codec is an interface to be implemented by any type that can encrypt and
// decrypt values
type Codec interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}
