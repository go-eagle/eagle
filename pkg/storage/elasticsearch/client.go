package elasticsearch

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	ErrEsNotFound = errors.New("es: not found")
)

type EsClient struct {
	client *es.Client
}

func NewESClient(hosts []string, username, password string) (client *EsClient, err error) {
	c, err := es.NewClient(es.Config{
		Addresses: hosts,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
			ResponseHeaderTimeout: 5 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &EsClient{client: c}, nil
}

func (es *EsClient) handleResponse(resp *esapi.Response) (map[string]interface{}, error) {
	var (
		r map[string]interface{}
	)
	if resp.StatusCode == 404 {
		return r, ErrEsNotFound
	}
	// Check response status
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing the response body: %s", err))
	}
	return r, nil
}
