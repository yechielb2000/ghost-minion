package communication

import (
	"bytes"
	"fmt"
	"ghostminion/config"
	"io"
	"net/http"
	"time"
)

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
)

func CreateRoute(serverConfig config.ServerConfig, route string) string {
	return fmt.Sprintf("https://%s:%d/%s", serverConfig.Address, serverConfig.Port, route)
}

func SendRequest(method HTTPMethod, url string, headers map[string]string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequest(string(method), url, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}
