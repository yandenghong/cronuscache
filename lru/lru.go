package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// an optional callback method when an entry is removed
	OnRemoved func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onRemoved func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnRemoved: onRemoved,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if item, ok := c.cache[key]; ok {
		c.ll.MoveToFront(item)
		kv := item.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Remove() {
	if item := c.ll.Back(); item != nil {
		c.ll.Remove(item)
		kv := item.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnRemoved != nil {
			c.OnRemoved(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// update
	if item, ok := c.cache[key]; ok {
		c.ll.MoveToFront(item)
		kv := item.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // add
		item := c.ll.PushFront(&entry{key, value})
		c.cache[key] = item
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.Remove()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
