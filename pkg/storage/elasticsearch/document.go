package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
)

const (
	defaultDocumentType = "_doc"
)

// CreateDocument create a document
func (es *EsClient) CreateDocument(ctx context.Context, indexName string, docId string,
	data map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}

	resp, err := es.client.Create(indexName, docId, &buf,
		es.client.Create.WithContext(ctx),
		es.client.Create.WithDocumentType(defaultDocumentType),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}

// GetDocument get a document
func (es *EsClient) GetDocument(ctx context.Context, indexName, docId string) (map[string]interface{}, error) {
	resp, err := es.client.Get(indexName, docId,
		es.client.Get.WithContext(ctx),
		es.client.Get.WithDocumentType(defaultDocumentType),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}

// IsExistDocument check a document is exist
func (es *EsClient) IsExistDocument(ctx context.Context, indexName, docId string) (bool, error) {
	resp, err := es.client.Exists(indexName, docId,
		es.client.Exists.WithContext(ctx),
		es.client.Exists.WithDocumentType(defaultDocumentType),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return false, ErrEsNotFound
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}

// UpdateDocument update a document
func (es *EsClient) UpdateDocument(ctx context.Context, indexName, docId string, data map[string]interface{}) (map[string]interface{}, error) {
	doc := map[string]interface{}{
		"doc": data,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(doc)
	if err != nil {
		return nil, err
	}
	resp, err := es.client.Update(indexName, docId, &buf,
		es.client.Update.WithContext(ctx),
		es.client.Update.WithDocumentType(defaultDocumentType),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}

// UpdateDocumentByQuery 根据搜索条件更新
func (es *EsClient) UpdateDocumentByQuery(ctx context.Context, indexName string, query, data map[string]interface{}) (map[string]interface{}, error) {
	doc := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query,
		},
		"script": map[string]interface{}{
			"inline": "ctx._source=params",
			"params": data,
			"lang":   "painless",
		},
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(doc)
	if err != nil {
		return nil, err
	}
	resp, err := es.client.UpdateByQuery(
		[]string{indexName},
		es.client.UpdateByQuery.WithContext(ctx),
		es.client.UpdateByQuery.WithDocumentType(defaultDocumentType),
		es.client.UpdateByQuery.WithBody(&buf),
		es.client.UpdateByQuery.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}

func (es *EsClient) DeleteDocument(ctx context.Context, indexName string, docId string) error {
	resp, err := es.client.Delete(indexName, docId,
		es.client.Delete.WithContext(ctx),
		es.client.Delete.WithDocumentType(defaultDocumentType),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = es.handleResponse(resp)
	return err
}

func (es *EsClient) DeleteDocumentByQuery(ctx context.Context, indexName string, query, data map[string]interface{}) (map[string]interface{}, error) {
	doc := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query,
		},
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(doc)
	if err != nil {
		return nil, err
	}
	resp, err := es.client.DeleteByQuery([]string{indexName}, &buf, es.client.DeleteByQuery.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return es.handleResponse(resp)
}
