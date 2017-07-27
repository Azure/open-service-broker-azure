package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func addAuthHeader(r *http.Request, username, password string) {
	usernameAndPassword := fmt.Sprintf("%s:%s", username, password)
	b64UsernameAndPassword := base64.StdEncoding.EncodeToString(
		[]byte(usernameAndPassword),
	)
	r.Header.Add(
		"Authorization",
		fmt.Sprintf("Basic %s", b64UsernameAndPassword),
	)
}
