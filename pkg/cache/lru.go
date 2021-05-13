// LRU cache
// wiki: https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
// impl1: https://dev.to/clavinjune/lru-cache-in-go-1cfk
// impl2: https://github.com/baidu/go-lib/blob/master/lru_cache/lru_cache.go
package cache

import "fmt"

type Node struct {
	Key   int
	Value int
	Prev  *Node
	Next  *Node
}

func NewNode(key, value int) *Node {
	return &Node{
		Key:   key,
		Value: value,
	}
}

type LRU struct {
	capacity int
	size     int
	data     map[int]*Node
	tail     *Node
	head     *Node
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		size:     0,
		data:     make(map[int]*Node),
	}
}

func (l *LRU) pushTail(n *Node) {
	if l.head == nil {
		l.head = n
		l.tail = n
		return
	}

	l.tail.Next = n
	n.Prev = l.tail
	l.tail = n
	l.tail.Next = nil
}

func (l *LRU) popHead() *Node {
	ret := l.head

	if l.head == l.tail {
		l.head = nil
	} else {
		l.head = l.head.Next
		l.head.Prev = nil
	}

	return ret
}

func (l *LRU) popTail() *Node {
	ret := l.tail

	if l.head == l.tail {
		l.head = nil
	} else {
		l.tail = l.tail.Prev
		l.tail.Next = nil
	}

	return ret
}

func (l *LRU) pop(n *Node) *Node {
	switch n {
	case l.head:
		return l.popHead()
	case l.tail:
		return l.popTail()
	}

	n.Next.Prev = n.Prev
	n.Prev.Next = n.Next
	return n
}

func (l *LRU) Set(key, value int) {
	// check if the key exists
	// if it exists, we need to remove it
	// then we append it to the queue
	// 4th rule (mark it as the most recently used)
	if val, isOk := l.data[key]; isOk {
		// this is the reason why we need to use popTail
		l.pop(val)
		l.size--
	}

	// 3rd rule
	if l.size >= l.capacity {
		n := l.popHead()
		delete(l.data, n.Key)
		l.size--
	}

	// push new data
	n := NewNode(key, value)
	l.data[key] = n
	l.pushTail(n)
	l.size++
}

func (l *LRU) Get(key int) int {
	val, isOk := l.data[key]

	if !isOk {
		return -1
	}

	// remove it
	l.pop(val)
	// then mark it as the most recently used
	l.pushTail(val)

	return val.Value
}

func (l *LRU) ShowQueue() {
	fmt.Printf("Least ")
	for n := l.head; n != l.tail; n = n.Next {
		fmt.Printf("%v -> ", n.Key)
	}

	fmt.Println(l.tail.Key, "Most")
}
