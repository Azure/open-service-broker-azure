package client

import (
	"bytes"
	"fmt"
	"net/http"
)

func getBaseURL(useSSL bool, host string, port int) string {
	proto := "http"
	if useSSL {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s:%d", proto, host, port)
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
