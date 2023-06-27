package elasticsearch

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	esClient *EsClient
	ctx      context.Context
)

const (
	indexName = "test-index"
)

func init() {
	var err error
	esClient, err = NewESClient([]string{"http://127.0.0.1:9200"}, "", "")
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.Background()
}

func TestNewESClient(t *testing.T) {
	t.Log(esClient.client.Info())
}

func TestNewESClientWithContext(t *testing.T) {
	t.Log(esClient.client.Info(esClient.client.Info.WithContext(context.Background())))
}

func TestCreateIndex(t *testing.T) {
	a := assert.New(t)
	bodyStr := `{
		  "settings": {
			"analysis": {
			  "analyzer": {
				"nickname_analyzer": {
				  "type": "custom",
				  "tokenizer": "ngram",
				  "char_filter": [
					"html_strip"
				  ],
				  "filter": [
					"lowercase",
					"trim",
					"unique"
				  ]
				}
			  },
			  "tokenizer": {
				"ngram": {
				  "type": "ngram",
				  "min_gram": 2,
				  "max_gram": 20
				}
			  }
			}
		  },
		  "mappings" : {
			  "properties" : {
				"id" : {
				  "index": true,
				  "type" : "keyword"
				},
				"age" : {
				  "index": true,
				  "type" : "keyword"
				},
				"nickname" : {
				  "index": true,
				  "type" : "text",
				  "analyzer": "partial"
				},
				"bio" : {
				  "index": true,
				  "type" : "keyword"
				}
			  }
	      }
		}`
	resp, err := esClient.CreateIndex(ctx, indexName, bodyStr)
	a.Nil(err)
	t.Log(resp)

	response, err := esClient.client.Indices.Get([]string{indexName})
	a.Nil(err)
	t.Log(response)
}

func TestGetIndex(t *testing.T) {
	a := assert.New(t)
	response, err := esClient.GetIndex(ctx, indexName)
	a.Nil(err)
	t.Log(response)
}

func TestDeleteIndex(t *testing.T) {
	a := assert.New(t)
	err := esClient.DeleteIndex(ctx, indexName)
	a.Nil(err)
}

func TestCreateDocument(t *testing.T) {
	a := assert.New(t)
	doc := map[string]interface{}{
		"id":       "1",
		"nickname": "Tom",
		"age":      "18",
		"bio":      "I am a singer",
	}
	response, err := esClient.CreateDocument(ctx, indexName, "1", doc)
	a.Nil(err)
	t.Log(response)
}

func TestGetDocument(t *testing.T) {
	a := assert.New(t)
	response, err := esClient.GetDocument(ctx, indexName, "1")
	a.Nil(err)
	t.Log(response)
}

func TestUpdateDocument(t *testing.T) {
	a := assert.New(t)
	data := map[string]interface{}{
		"bio": "I am a teacher",
	}
	response, err := esClient.UpdateDocument(ctx, indexName, "1", data)
	a.Nil(err)
	t.Log(response)
}

func TestUpdateDocumentByQuery(t *testing.T) {
	a := assert.New(t)
	query := map[string]interface{}{
		"nickname": "Tom",
	}
	data := map[string]interface{}{
		"bio": "I am a music teacher",
	}
	response, err := esClient.UpdateDocumentByQuery(ctx, indexName, query, data)
	a.Nil(err)
	t.Log(response)
}

func TestDeleteDocument(t *testing.T) {
	a := assert.New(t)
	err := esClient.DeleteDocument(ctx, indexName, "1")
	a.Nil(err)
}

// 单个字段的搜索
func TestSearch(t *testing.T) {
	a := assert.New(t)
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"match": map[string]interface{}{
						"nickname": "Tom",
					},
				},
			},
		},
	}
	sort := map[string]interface{}{
		"_score": map[string]interface{}{
			"order": "desc",
		},
	}
	response, total, err := esClient.Search(ctx, indexName, query, sort, 0, 10)
	a.Nil(err)
	t.Log(total)
	t.Log(response)
}

// 多字段搜索
func TestSearchByMultiField(t *testing.T) {
	a := assert.New(t)
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"match": map[string]interface{}{
						"nickname": "Tom",
					},
				},
				{
					"match": map[string]interface{}{
						"age": "18",
					},
				},
			},
		},
	}
	sort := map[string]interface{}{
		"_score": map[string]interface{}{
			"order": "desc",
		},
	}
	response, total, err := esClient.Search(ctx, indexName, query, sort, 0, 10)
	a.Nil(err)
	t.Log(total)
	t.Log(response)
}
