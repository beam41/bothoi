package http_util

import (
	"io"
	"log"
	"net/http"
)

func Get(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
