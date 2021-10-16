package cache

// Getter .
type Getter interface {
	Get(key string) interface{}
}

// GetFunc .
type GetFunc func(key string) interface{}

// Get .
func (f GetFunc) Get(key string) interface{} {
	return f(key)
}
