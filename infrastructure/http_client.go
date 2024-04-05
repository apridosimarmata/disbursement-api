package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Get(ctx context.Context, path string, queryParams map[string]string, result interface{}) (err error)
	Post(ctx context.Context, path string, body map[string]interface{}, result interface{}) (err error)
}

type httpClient struct {
	serverBaseUrl string
	client        http.Client
}

func NewHTTPClient(serverBaseUrl string) HTTPClient {
	return &httpClient{
		client:        http.Client{},
		serverBaseUrl: serverBaseUrl,
	}
}

func (client *httpClient) Get(ctx context.Context, path string, queryParams map[string]string, result interface{}) (err error) {
	reqUrl := client.serverBaseUrl + path
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		Log(fmt.Sprintf("%s - http.NewRequest @ httpClient.Get", err.Error()))
		return err
	}

	if queryParams != nil {
		q := request.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	response, err := client.client.Do(request)
	if err != nil {
		Log(fmt.Sprintf("%s - client.client.Do @ httpClient.Get", err.Error()))
		return err
	}
	defer response.Body.Close()

	if err = client.validateResponseBasedOnStatusCode(response.StatusCode); err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		Log(fmt.Sprintf("%s - json.NewDecoder(response.Body).Decode @ httpClient.Get", err.Error()))
		return err
	}

	return nil
}

func (client *httpClient) Post(ctx context.Context, path string, body map[string]interface{}, result interface{}) (err error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		Log(fmt.Sprintf("%s - json.Marshal @ httpClient.Post", err.Error()))
		return err
	}

	request, err := http.NewRequest("POST", client.serverBaseUrl+path, bytes.NewBuffer(reqBody))
	if err != nil {
		Log(fmt.Sprintf("%s - http.NewRequest @ httpClient.Post", err.Error()))
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.client.Do(request)
	if err != nil {
		Log(fmt.Sprintf("%s - client.client.Do @ httpClient.Post", err.Error()))
		return err
	}
	defer response.Body.Close()

	if err = client.validateResponseBasedOnStatusCode(response.StatusCode); err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		Log(fmt.Sprintf("%s - json.NewDecoder(response.Body).Decode @ httpClient.Post", err.Error()))
		return err
	}

	return nil
}

func (client *httpClient) validateResponseBasedOnStatusCode(code int) (err error) {
	if code == http.StatusOK {
		return nil
	}

	switch {
	case 400 >= code && code <= 499:
		return errors.New("client error")
	case 500 >= code && code <= 599:
		return errors.New("server error")
	}

	return errors.New("expecting OK status code")
}
