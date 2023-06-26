package elasticsearch

import (
	"context"
	"strings"
)

// CreateIndex creates an index with the given name and body.
func (es *EsClient) CreateIndex(ctx context.Context, indexName string, body string) (map[string]interface{}, error) {
	resp, err := es.client.Indices.Create(indexName,
		es.client.Indices.Create.WithContext(ctx),
		es.client.Indices.Create.WithBody(strings.NewReader(body)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}

// GetIndex gets an index with the given name.
func (es *EsClient) GetIndex(ctx context.Context, indexName string) (map[string]interface{}, error) {
	resp, err := es.client.Indices.Get([]string{indexName}, es.client.Indices.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}

// DeleteIndex deletes an index with the given name.
func (es *EsClient) DeleteIndex(ctx context.Context, indexName string) error {
	resp, err := es.client.Indices.Delete([]string{indexName}, es.client.Indices.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = es.handleResponse(resp)
	return err
}
