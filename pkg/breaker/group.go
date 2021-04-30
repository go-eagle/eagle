package breaker

import "sync"

// Group represents a class of CircuitBreaker and forms a namespace in which
// units of CircuitBreaker.
type Group struct {
	mu   sync.RWMutex
	brks map[string]Breaker
	conf *Config
}

// NewGroup new a breaker group container, if conf nil use default conf.
func NewGroup(conf *Config) *Group {
	if conf == nil {
		_mu.RLock()
		conf = _conf
		_mu.RUnlock()
	} else {
		conf.fix()
	}
	return &Group{
		conf: conf,
		brks: make(map[string]Breaker),
	}
}

// Get get a breaker by a specified key, if breaker not exists then make a new one.
func (g *Group) Get(key string) Breaker {
	g.mu.RLock()
	brk, ok := g.brks[key]
	conf := g.conf
	g.mu.RUnlock()
	if ok {
		return brk
	}
	// NOTE here may new multi breaker for rarely case, let gc drop it.
	brk = newBreaker(conf)
	g.mu.Lock()
	if _, ok = g.brks[key]; !ok {
		g.brks[key] = brk
	}
	g.mu.Unlock()
	return brk
}

// Reload reload the group by specified config, this may let all inner breaker
// reset to a new one.
func (g *Group) Reload(conf *Config) {
	if conf == nil {
		return
	}
	conf.fix()
	g.mu.Lock()
	g.conf = conf
	g.brks = make(map[string]Breaker, len(g.brks))
	g.mu.Unlock()
}

// Go runs your function while tracking the breaker state of group.
func (g *Group) Go(name string, run, fallback func() error) error {
	breaker := g.Get(name)
	if err := breaker.Allow(); err != nil {
		return fallback()
	}
	return run()
}
