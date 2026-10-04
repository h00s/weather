package services

import (
	"context"
	"sync"
	"time"
)

// failureBackoff is how long a failed upstream is left alone: callers meanwhile get
// the last good value (or the error) at once, rather than queuing for another try.
const failureBackoff = 30 * time.Second

// staleCache holds one upstream value for ttl. When a refresh fails it keeps
// serving the last good value, so the display shows slightly old data instead
// of none; it returns the error only while nothing has been fetched yet.
type staleCache[T any] struct {
	ttl time.Duration
	now func() time.Time

	mu        sync.Mutex
	value     T
	fetchedAt time.Time
	failedAt  time.Time // the last failed refresh; zero after a success
	failure   error
}

func newStaleCache[T any](ttl time.Duration) *staleCache[T] {
	return &staleCache[T]{ttl: ttl, now: time.Now}
}

// Get returns the cached value and when it was fetched, calling fetch when the
// value is older than ttl. The lock is held across fetch, so concurrent callers
// share one upstream request; after a failure, callers within failureBackoff
// share its outcome too.
func (c *staleCache[T]) Get(ctx context.Context, fetch func(context.Context) (T, error)) (T, time.Time, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.fetchedAt.IsZero() && c.now().Sub(c.fetchedAt) < c.ttl {
		return c.value, c.fetchedAt, nil
	}
	if !c.failedAt.IsZero() && c.now().Sub(c.failedAt) < failureBackoff {
		return c.stale(c.failure)
	}

	value, err := fetch(ctx)
	if err != nil {
		if ctx.Err() == nil { // a caller that gave up is no verdict on the upstream
			c.failedAt, c.failure = c.now(), err
		}
		return c.stale(err)
	}

	c.value, c.fetchedAt, c.failedAt, c.failure = value, c.now(), time.Time{}, nil
	return c.value, c.fetchedAt, nil
}

// stale is the last good value, or err when there is none.
func (c *staleCache[T]) stale(err error) (T, time.Time, error) {
	if c.fetchedAt.IsZero() {
		var zero T
		return zero, time.Time{}, err
	}
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
