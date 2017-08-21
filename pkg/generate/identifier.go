package generate

const (
	identifierLength = 10
	identifierChars  = lowerAlphaChars + numberChars
)

// NewIdentifier generates an identifier suitable for use as a username,
// role name, database name for various database systems-- including, at least,
// PostgresSQL and MySQL.
func NewIdentifier() string {
	seededRandMutex.Lock()
	defer seededRandMutex.Unlock()
	b := make([]byte, identifierLength)
	// The first character of an identifier MUST be a lowercase alpha
	b[0] = lowerAlphaChars[seededRand.Intn(len(lowerAlphaChars))]
	// The rest can be lowercase alphas or numerals
	for i := 1; i < identifierLength; i++ {
		b[i] = identifierChars[seededRand.Intn(len(identifierChars))]
	}
	return string(b)
}
