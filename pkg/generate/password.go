package generate

const (
	passwordLength = 16
	passwordChars  = lowerAlphaChars + upperAlphaChars + numberChars
)

// NewPassword generates a strong, random password
func NewPassword() string {
	b := make([]byte, passwordLength)
	// Passwords need to include at least one character from each of the three
	// groups. To ensure that, we'll fill each of the first three []byte elements
	// with a random character from a specific group.
	b[0] = lowerAlphaChars[seededRand.Intn(len(lowerAlphaChars))]
	b[1] = upperAlphaChars[seededRand.Intn(len(upperAlphaChars))]
	b[2] = numberChars[seededRand.Intn(len(numberChars))]
	// The remainder of the characters can be completely random and drawn from
	// all three character groups.
	for i := 3; i < passwordLength; i++ {
		b[i] = passwordChars[seededRand.Intn(len(passwordChars))]
	}
	// For good measure, shuffle the elements of the entire []byte so that
	// the 0 character isn't predicatably lowercase, etc...
	for i := range b {
		j := seededRand.Intn(len(b))
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
