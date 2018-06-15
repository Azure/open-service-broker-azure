package client

import (
	"bytes"
	"fmt"
	"net/http"
)

func getBaseURL(host string, port int) string {
	return fmt.Sprintf("http://%s:%d", host, port)
}

func newRequest(
	method string,
	url string,
	username string,
	password string,
	body []byte,
) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error building request: %s", err)
	}
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}
	req.Header.Add("X-Broker-API-Version", "2.13")
	return req, nil
}
