package services

import (
	"context"
	"sync"
	"time"
)

// staleCache holds one upstream value for ttl. When a refresh fails it keeps
// serving the last good value, so the display shows slightly old data instead
// of none; it returns the error only while nothing has been fetched yet.
type staleCache[T any] struct {
	ttl time.Duration
	now func() time.Time

	mu        sync.Mutex
	value     T
	fetchedAt time.Time
}

func newStaleCache[T any](ttl time.Duration) *staleCache[T] {
	return &staleCache[T]{ttl: ttl, now: time.Now}
}

// Get returns the cached value and when it was fetched, calling fetch when the
// value is older than ttl. The lock is held across fetch, so concurrent callers
// share one upstream request.
func (c *staleCache[T]) Get(ctx context.Context, fetch func(context.Context) (T, error)) (T, time.Time, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.fetchedAt.IsZero() && c.now().Sub(c.fetchedAt) < c.ttl {
		return c.value, c.fetchedAt, nil
	}

	value, err := fetch(ctx)
	if err != nil {
		if c.fetchedAt.IsZero() {
			var zero T
			return zero, time.Time{}, err
		}
		return c.value, c.fetchedAt, nil
	}

	c.value, c.fetchedAt = value, c.now()
	return c.value, c.fetchedAt, nil
}

// keyedCache keeps a staleCache per key, for upstreams asked per location or
// per query. It holds at most max keys: adding one more drops the least
// recently used. Callers with different keys never wait on each other.
type keyedCache[K comparable, T any] struct {
	ttl time.Duration
	max int
	now func() time.Time

	mu      sync.Mutex
	entries map[K]*keyedEntry[T]
}

type keyedEntry[T any] struct {
	cache    *staleCache[T]
	lastUsed time.Time
}

func newKeyedCache[K comparable, T any](ttl time.Duration, max int) *keyedCache[K, T] {
	return &keyedCache[K, T]{ttl: ttl, max: max, now: time.Now, entries: map[K]*keyedEntry[T]{}}
}

// Get returns key's value as staleCache.Get does.
func (c *keyedCache[K, T]) Get(ctx context.Context, key K, fetch func(context.Context) (T, error)) (T, time.Time, error) {
	return c.entry(key).Get(ctx, fetch)
}

func (c *keyedCache[K, T]) entry(key K) *staleCache[T] {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		if len(c.entries) >= c.max {
			c.evictLeastRecentlyUsed()
		}
		e = &keyedEntry[T]{cache: &staleCache[T]{ttl: c.ttl, now: c.now}}
		c.entries[key] = e
	}
	e.lastUsed = c.now()
	return e.cache
}

func (c *keyedCache[K, T]) evictLeastRecentlyUsed() {
	var oldest K
	var oldestAt time.Time
	first := true
	for k, e := range c.entries {
		if first || e.lastUsed.Before(oldestAt) {
			oldest, oldestAt, first = k, e.lastUsed, false
		}
	}
	delete(c.entries, oldest)
}
