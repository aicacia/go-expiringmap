package expiringmap

import (
	"time"

	"github.com/aicacia/go-cmap"
)

type expiringMapVal[V any] struct {
	val V
	ttl time.Time
}

type ExpiringMap[K, V any] struct {
	items cmap.CMap[K, expiringMapVal[V]]
}

func New[K, V any]() ExpiringMap[K, V] {
	return ExpiringMap[K, V]{
		items: cmap.New[K, expiringMapVal[V]](),
	}
}

func (m *ExpiringMap[K, V]) SetIfAbsent(key K, value V, ttl time.Time) bool {
	return m.items.SetIfAbsent(key, expiringMapVal[V]{value, ttl})
}

func (m *ExpiringMap[K, V]) Set(key K, value V, ttl time.Time) bool {
	return m.items.Set(key, expiringMapVal[V]{value, ttl})
}

func (m *ExpiringMap[K, V]) GetOrSet(key K, value V, ttl time.Time) V {
	newItem := expiringMapVal[V]{value, ttl}
	item, _ := m.items.LoadOrStore(key, newItem)
	if item.ttl.Before(time.Now()) {
		m.items.Set(key, newItem)
		return value
	}
	return item.val
}

func (m *ExpiringMap[K, V]) Has(key K) bool {
	if item, ok := m.items.Get(key); ok {
		if item.ttl.Before(time.Now()) {
			m.items.Delete(key)
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func (m *ExpiringMap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *ExpiringMap[K, V]) Get(key K) (V, bool) {
	if item, ok := m.items.Get(key); ok {
		if item.ttl.Before(time.Now()) {
			m.items.Delete(key)
		} else {
			return item.val, true
		}
	}
	return *new(V), false
}

func (m *ExpiringMap[K, V]) Delete(key K) bool {
	return m.items.Delete(key)
}

func (m *ExpiringMap[K, V]) Remove(key K) bool {
	return m.items.Remove(key)
}

func (m *ExpiringMap[K, V]) Range(f func(key K, value V) bool) {
	now := time.Now()
	m.items.Range(func(key K, value expiringMapVal[V]) bool {
		if value.ttl.Before(now) {
			m.items.Delete(key)
			return true
		} else {
			return f(key, value.val)
		}
	})
}

func (m *ExpiringMap[K, V]) Iter() chan cmap.Entry[K, V] {
	ch := make(chan cmap.Entry[K, V])
	go func() {
		m.Range(func(key K, value V) bool {
			ch <- cmap.Entry[K, V]{
				Key: key,
				Val: value,
			}
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *ExpiringMap[K, V]) Keys() chan K {
	ch := make(chan K)
	go func() {
		m.Range(func(key K, _ V) bool {
			ch <- key
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *ExpiringMap[K, V]) Values() chan V {
	ch := make(chan V)
	go func() {
		m.Range(func(_ K, value V) bool {
			ch <- value
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *ExpiringMap[K, V]) Len() int {
	count := 0
	m.Range(func(_ K, _ V) bool {
		count += 1
		return true
	})
	return count
}

func (m *ExpiringMap[K, V]) Clear() {
	m.items.Clear()
}
