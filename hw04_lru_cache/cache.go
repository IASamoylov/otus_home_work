package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if v, ok := l.items[key]; ok {
		l.queue.MoveToFront(v)
		cacheItem := v.Value.(cacheItem)
		cacheItem.value = value
		v.Value = cacheItem
		return true
	}

	if l.capacity == l.queue.Len() {
		v := l.queue.Back()
		cacheItem := v.Value.(cacheItem)
		delete(l.items, cacheItem.key)
		l.queue.Remove(v)
	}

	l.items[key] = l.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if v, ok := l.items[key]; ok {
		l.queue.MoveToFront(v)
		cacheItem := v.Value.(cacheItem)
		return cacheItem.value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue.Clear()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type cacheItem struct {
	key   Key
	value interface{}
}
