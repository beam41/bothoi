package http_util

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// post json
func PostJson(url string, body any) ([]byte, error) {
	header := map[string]string{}
	return PostJsonH(url, body, header)
}

// post json with header
func PostJsonH(url string, body any, header map[string]string) ([]byte, error) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(body)
	if err != nil {
		return nil, err
	}
	header["Content-Type"] = "application/json"
	return Post(url, payload, header)
}

func Post(url string, body *bytes.Buffer, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	for h, v := range header {
		req.Header.Add(h, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("Post StatusCode", res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		return nil, err
	}
	return resBody, nil
}
