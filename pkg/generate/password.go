package generate

const (
	passwordLength = 16
	passwordChars  = lowerAlphaChars + upperAlphaChars + numberChars
)

// NewPassword generates a strong, random password
func NewPassword() string {
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = passwordChars[seededRand.Intn(len(passwordChars))]
	}
	return string(b)
}
