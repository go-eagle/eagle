package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// Search documents by query and sort
func (es *EsClient) Search(ctx context.Context, indexName string, query, sort map[string]interface{}, offset,
	limit int) (ret []map[string]interface{}, total int, err error) {
	var (
		r map[string]interface{}
	)
	query = map[string]interface{}{
		"query": query,
		"sort":  sort,
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return nil, 0, err
	}
	resp, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(indexName),
		es.client.Search.WithDocumentType(defaultDocumentType),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackScores(true),
		es.client.Search.WithFrom(offset),
		es.client.Search.WithSize(limit),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	fmt.Println("===resp===", resp)
	r, err = es.handleResponse(resp)
	if err != nil {
		return ret, 0, err
	}
	total = int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		retTmp := make(map[string]interface{})
		retTmp["_id"] = hit.(map[string]interface{})["_id"]
		source := hit.(map[string]interface{})["_source"]
		// _source is biz data
		retTmp["data"] = source
		ret = append(ret, retTmp)
	}
	return ret, total, nil
}

// SearchByAggregations documents by body with json format
func (es *EsClient) SearchByAggregations(ctx context.Context, indexName string,
	jsonBody string) (ret map[string]interface{}, err error) {
	var (
		r map[string]interface{}
	)

	resp, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(indexName),
		es.client.Search.WithBody(bytes.NewBufferString(jsonBody)),
		es.client.Search.WithTrackScores(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err = es.handleResponse(resp)
	if err != nil {
		return ret, err
	}

	// aggregations need to custom handle because of value of aggregations is not fixed
	return r["aggregations"].(map[string]interface{}), nil
}
