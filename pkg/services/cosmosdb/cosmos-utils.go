package cosmosdb

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func generateAuthToken(verb, resource, id, date, key string) (string, error) {
	payload := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		strings.ToLower(verb),
		strings.ToLower(resource),
		id,
		strings.ToLower(date),
		"",
	)

	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	hmac := hmac.New(sha256.New, decodedKey)
	_, err := hmac.Write([]byte(payload))
	if err != nil {
		return "", err
	}
	b := hmac.Sum(nil)
	authHash := base64.StdEncoding.EncodeToString(b)
	authHeader := url.QueryEscape("type=master&ver=1.0&sig=" + authHash)
	return authHeader, nil
}

func createRequest(
	accountName string,
	method string,
	resourceType string,
	resourceID string,
	key string,
	body interface{},
) (*http.Request, error) {
	path := fmt.Sprintf("%s/%s", resourceType, resourceID)
	url := fmt.Sprintf("https://%s.documents.azure.com/%s", accountName, path)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest(method, url, buf)
	dateStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	authHeader, err := generateAuthToken(
		method,
		resourceType,
		resourceID,
		dateStr,
		key,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Ms-Date", dateStr)
	req.Header.Add("X-Ms-version", "2017-02-22")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	return req, nil
}

func createDatabase(
	accountName string,
	id string,
	key string,
) error {
	request := &databaseCreationRequest{
		ID: id,
	}
	req, err := createRequest(accountName, "POST", "dbs", "", key, request)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("error creating database")
	}
	return nil
}
