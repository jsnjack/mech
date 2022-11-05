package cmd

import (
	"fmt"
	"io"
	"net/http"
)

func makeAPIRequest(method string, url string, payload io.Reader) (respBody []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("x-cns-security-token", buildSecurityToken())
	req.Header.Add("Content-Type", "application/json")
	if rootVerbose {
		fmt.Printf("requesting %s %s ...\n", method, url)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return body, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return body, nil
}
