package utils

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

func HTTPRequest(ctx context.Context, endpoint string, method string, payload []byte, contentType string, headers map[string]interface{}) (*http.Response, error) {
	log.Printf("Initializing Request: %v\n", endpoint)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	for headerString, header := range headers {
		req.Header.Set(headerString, header.(string))
	}

	client := &http.Client{
		Timeout: time.Second * 60,
	}

	req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println("Ending Request")

	return res, nil
}
