package postgresql

// generatePassword generates a random password
func generatePassword() string {
	seededRandMutex.Lock()
	defer seededRandMutex.Unlock()
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = passwordChars[seededRand.Intn(len(passwordChars))]
	}
	return string(b)
}
