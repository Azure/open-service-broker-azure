package client

import (
	"crypto/tls"
	"net/http"
)

func getHTTPClient(skipCertValidation bool) *http.Client {
	if !skipCertValidation {
		return &http.Client{}
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
