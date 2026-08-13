package lrucache

import "sync"

type LRUCache struct {
	capacity int
	items    map[string]*Node
	mu       sync.RWMutex
	head     *Node
	tail     *Node
}

type Node struct {
	Key   string
	Value int
	Next  *Node
	Prev  *Node
}

func NewCache(capacity int) *LRUCache {
	items := make(map[string]*Node, capacity)
	return &LRUCache{
		capacity: capacity,
		items:    items,
		mu:       sync.RWMutex{},
	}
}

func (c *LRUCache) Set(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	node, ok := c.items[key]
	if ok {
		node.Value = value
		c.remove(node)
		c.addToFront(node)
		return
	}

	newNode := &Node{
		Key:   key,
		Value: value,
	}

	c.items[key] = newNode
	c.addToFront(newNode)
	if len(c.items) > c.capacity {
		delete(c.items, c.tail.Key)
		c.remove(c.tail)
	}
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *LRUCache) Get(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	node, ok := c.items[key]
	if !ok {
		return -1, false
	}

	c.remove(node)
	c.addToFront(node)
	return node.Value, true
}

func (c *LRUCache) addToFront(node *Node) {
	node.Next = c.head
	if c.head != nil {
		c.head.Prev = node
	}

	c.head = node
	if c.tail == nil {
		c.tail = node
	}
}

func (c *LRUCache) remove(node *Node) {
	if node == c.head && node == c.tail {
		c.head, c.tail = nil, nil
		return
	}

	if node == c.head {
		next := c.head.Next
		c.head.Next = nil
		next.Prev = nil
		c.head = next
		return
	}

	if node == c.tail {
		prev := c.tail.Prev
		c.tail.Prev = nil
		prev.Next = nil
		c.tail = prev
		return
	}

	next := node.Next
	prev := node.Prev
	node.Next = nil
	node.Prev = nil
	prev.Next = next
	next.Prev = prev
}
