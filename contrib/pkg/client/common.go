package client

import "fmt"

func getBaseURL(host string, port int) string {
	return fmt.Sprintf("http://%s:%d", host, port)
}
