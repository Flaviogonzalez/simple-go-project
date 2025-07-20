package helpers

import (
	"bytes"
	"io"
	"net/http"
)

func SendRequest(url string, method string, data []byte) (r []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}
