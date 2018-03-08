package api

type errBasicAuthUsernameNotSpecified struct{}

func (e *errBasicAuthUsernameNotSpecified) Error() string {
	return "neither environment variable BASIC_AUTH_USERNAME nor " +
		"SECURITY_USER_NAME was specified; please specify one or the other"
}

type errBasicAuthPasswordNotSpecified struct{}

func (e *errBasicAuthPasswordNotSpecified) Error() string {
	return "neither environment variable BASIC_AUTH_PASSWORD nor " +
		"SECURITY_USER_PASSWORD was specified; please specify one or the other"
}
