package headers

var responseMissingAPIVersion = []byte(
	`{ "error": "MissingAPIVersion", "description": "The request did not ` +
		`include the X-Broker-API-Version header"}`,
)

var responseAPIVersionIncorrect = []byte(
	`{ "error": "APIVersionIncorrect", "description": "X-Broker-API-Verson ` +
		`header includes an incompatible version"}`,
)
