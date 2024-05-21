package lrucache

import (
	"container/list"
	"errors"
	"sync"
)

type entry struct {
	key   string
	value interface{}
}

var ErrCreate = errors.New("cap must be a positive num")

type LruCache struct {
	cap   int
	ll    *list.List
	mu    sync.RWMutex
	cache map[string]*list.Element
}

func NewLruCache(cap int) (*LruCache, error) {
	if cap <= 0 {
		return nil, ErrCreate
	}

	return &LruCache{
		cap:   cap,
		ll:    list.New(),
		cache: make(map[string]*list.Element),
	}, nil
}

func (l *LruCache) removeOldest() {

	ele := l.ll.Back()
	if ele != nil {
		k := ele.Value.(*entry).key
		delete(l.cache, k)
		l.ll.Remove(ele)
	}

}

func (l *LruCache) Set(key string, value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cache == nil {
		l.ll = list.New()
		l.cache = make(map[string]*list.Element)
	}

	if e, ok := l.cache[key]; ok {
		e.Value.(*entry).value = value
		l.ll.MoveToFront(e)
		return
	}

	ele := l.ll.PushFront(&entry{key, value})
	l.cache[key] = ele

	if l.cap != 0 && l.ll.Len() > l.cap {
		l.removeOldest()
	}
}

func (l *LruCache) Get(key string) (value interface{}, ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.cache == nil {
		return
	}

	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}

	return
}

func (l *LruCache) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.cache == nil {
		return 0
	}

	return l.ll.Len()
}

func (l *LruCache) remove(e *list.Element) {
	k := e.Value.(*entry).key
	delete(l.cache, k)
	l.ll.Remove(e)
}

func (l *LruCache) Remove(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cache == nil {
		return
	}

	if ele, ok := l.cache[key]; ok {
		l.remove(ele)
	}

}
