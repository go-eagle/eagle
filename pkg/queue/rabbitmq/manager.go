package rabbitmq

import (
	"fmt"
	"sync"

	"github.com/go-eagle/eagle/pkg/log"
)

var (
	DefaultManager *Manager
)

type Manager struct {
	opts map[string]*Config

	cmu *sync.RWMutex
	pmu *sync.RWMutex

	consumers  map[string]*Consumer
	publishers map[string]*Producer
}

func NewManager(opts map[string]*Config) *Manager {
	return &Manager{
		opts:       opts,
		cmu:        &sync.RWMutex{},
		pmu:        &sync.RWMutex{},
		consumers:  make(map[string]*Consumer, 0),
		publishers: make(map[string]*Producer, 0),
	}
}

func (c *Manager) GetProducer(name string) (*Producer, error) {
	c.pmu.Lock()
	defer c.pmu.Unlock()

	if _, ok := c.publishers[name]; ok {
		return c.publishers[name], nil
	}
	if opt, ok := c.opts[name]; ok {
		publisher, err := NewProducer(opt, log.GetLogger())
		if err != nil {
			return nil, err
		}
		c.publishers[name] = publisher
		return publisher, nil
	}
	return nil, fmt.Errorf("rabbitmq: GetPublisher error, config %s not found", name)
}

func (c *Manager) GetConsumer(name string) (*Consumer, error) {
	c.cmu.Lock()
	defer c.cmu.Unlock()

	if _, ok := c.consumers[name]; ok {
		return c.consumers[name], nil
	}
	if opt, isOk := c.opts[name]; isOk {
		consumer, err := NewConsumer(opt, log.GetLogger())
		if err != nil {
			return nil, err
		}
		c.consumers[name] = consumer
		return consumer, nil
	}
	return nil, fmt.Errorf("rabbitmq: GetConsumer error, config %s not found", name)
}

func (c *Manager) Close() error {
	for _, consumer := range c.consumers {
		_ = consumer.Close()
	}
	for _, publisher := range c.publishers {
		_ = publisher.Close()
	}
	return nil
}
