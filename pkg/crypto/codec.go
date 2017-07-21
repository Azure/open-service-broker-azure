package crypto

// Codec is an interface to be implemented by any type that can encrypt and
// decrypt values
type Codec interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}
