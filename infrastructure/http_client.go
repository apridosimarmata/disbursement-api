package infrastructure

import (
	"bytes"
	"context"
	"disbursement/domain/common/response"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(ctx context.Context, path string, queryParams map[string]string, result interface{}) (err error)
	Post(ctx context.Context, path string, body interface{}, result interface{}) (err error)
}

type httpClient struct {
	serverBaseUrl string
	client        http.Client
}

func NewHTTPClient(serverBaseUrl string) HTTPClient {
	return &httpClient{
		client: http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: time.Second * 6,
			},
		},
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

func (client *httpClient) Post(ctx context.Context, path string, body interface{}, result interface{}) (err error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		Log(fmt.Sprintf("%s - json.Marshal @ httpClient.Post", err.Error()))
		return err
	}

	reqUrl := client.serverBaseUrl + path
	request, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(reqBody))
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
	case code == 404:
		return errors.New(response.ERROR_NOT_FOUND)
	case code >= 400:
		return errors.New("client error")
	case code >= 500 && code <= 599:
		return errors.New("server error")
	}

	return errors.New("expecting OK status code")
}
