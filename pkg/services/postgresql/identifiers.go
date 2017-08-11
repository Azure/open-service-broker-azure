package postgresql

// generateIdentifier generates a valid PostgreSQL identifier. These can be
// used as database names and role names.
func generateIdentifier() string {
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
