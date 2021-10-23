package etcdclient

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3/mock/mockserver"
)

func startMockServer() {
	ms, err := mockserver.StartMockServers(1)
	if err != nil {
		log.Fatal(err)
	}

	if err := ms.StartAt(0); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	go startMockServer()
}

func Test_newClient(t *testing.T) {
	config := &Config{}
	config.Endpoints = []string{"localhost:0"}
	config.TTL = 5
	_, err := newClient(config)
	assert.Nil(t, err)
}
