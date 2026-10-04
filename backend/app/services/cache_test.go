package services

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time          { return c.now }
func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

func newTestCache(ttl time.Duration) (*staleCache[int], *fakeClock) {
	clock := &fakeClock{now: time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)}
	return &staleCache[int]{ttl: ttl, now: clock.Now}, clock
}

func counter(values ...int) (func(context.Context) (int, error), *int) {
	calls := 0
	return func(context.Context) (int, error) {
		v := values[min(calls, len(values)-1)]
		calls++
		return v, nil
	}, &calls
}

func TestStaleCacheServesFreshValueWithoutRefetching(t *testing.T) {
	c, clock := newTestCache(10 * time.Minute)
	fetch, calls := counter(1, 2)

	c.Get(t.Context(), fetch)
	clock.Advance(9 * time.Minute)
	v, fetchedAt, err := c.Get(t.Context(), fetch)

	if err != nil || v != 1 || *calls != 1 {
		t.Errorf("got %d, %v after %d fetches; want the cached 1 after 1 fetch", v, err, *calls)
	}
	if want := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC); !fetchedAt.Equal(want) {
		t.Errorf("fetchedAt = %v, want %v", fetchedAt, want)
	}
}

func TestStaleCacheRefetchesAfterTTL(t *testing.T) {
	c, clock := newTestCache(10 * time.Minute)
	fetch, calls := counter(1, 2)

	c.Get(t.Context(), fetch)
	clock.Advance(11 * time.Minute)
	v, fetchedAt, err := c.Get(t.Context(), fetch)

	if err != nil || v != 2 || *calls != 2 {
		t.Errorf("got %d, %v after %d fetches; want 2 after 2 fetches", v, err, *calls)
	}
	if !fetchedAt.Equal(clock.Now()) {
		t.Errorf("fetchedAt = %v, want %v", fetchedAt, clock.Now())
	}
}

func TestStaleCacheServesLastValueWhenRefreshFails(t *testing.T) {
	c, clock := newTestCache(10 * time.Minute)
	fetch, _ := counter(1)
	c.Get(t.Context(), fetch)
	firstFetch := clock.Now()
	clock.Advance(time.Hour)

	v, fetchedAt, err := c.Get(t.Context(), func(context.Context) (int, error) { return 0, errors.New("upstream down") })

	if err != nil || v != 1 {
		t.Errorf("got %d, %v; want the stale 1 and no error", v, err)
	}
	if !fetchedAt.Equal(firstFetch) {
		t.Errorf("fetchedAt = %v, want the original %v", fetchedAt, firstFetch)
	}
}

func TestStaleCacheReturnsErrorWhenNothingIsCached(t *testing.T) {
	c, _ := newTestCache(10 * time.Minute)
	boom := errors.New("upstream down")

	_, _, err := c.Get(t.Context(), func(context.Context) (int, error) { return 0, boom })

	if !errors.Is(err, boom) {
		t.Errorf("err = %v, want %v", err, boom)
	}
}

func newTestKeyedCache(ttl time.Duration, max int) (*keyedCache[string, int], *fakeClock) {
	clock := &fakeClock{now: time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)}
	c := newKeyedCache[string, int](ttl, max)
	c.now = clock.Now
	return c, clock
}

func TestKeyedCacheFetchesEachKeyOnce(t *testing.T) {
	c, _ := newTestKeyedCache(time.Minute, 10)
	fetchA, callsA := counter(1)
	fetchB, callsB := counter(2)

	for range 3 {
		if v, _, err := c.Get(t.Context(), "a", fetchA); v != 1 || err != nil {
			t.Fatalf("Get(a) = %d, %v", v, err)
		}
		if v, _, err := c.Get(t.Context(), "b", fetchB); v != 2 || err != nil {
			t.Fatalf("Get(b) = %d, %v", v, err)
		}
	}
	if *callsA != 1 || *callsB != 1 {
		t.Errorf("fetches: a %d, b %d; want one each", *callsA, *callsB)
	}
}

func TestKeyedCacheRefetchesAKeyAfterTTL(t *testing.T) {
	c, clock := newTestKeyedCache(time.Minute, 10)
	fetch, calls := counter(1, 2)

	c.Get(t.Context(), "a", fetch)
	clock.Advance(2 * time.Minute)
	if v, _, _ := c.Get(t.Context(), "a", fetch); v != 2 || *calls != 2 {
		t.Errorf("after TTL: value %d after %d fetches, want 2 after 2", v, *calls)
	}
}

func TestKeyedCacheDropsTheLeastRecentlyUsedKeyWhenFull(t *testing.T) {
	c, clock := newTestKeyedCache(time.Hour, 2)
	fetch, calls := counter(7)
	get := func(key string) {
		clock.Advance(time.Second)
		c.Get(t.Context(), key, fetch)
	}

	get("a")
	get("b")
	get("a") // b is now the least recently used
	get("c") // full: drops b
	before := *calls
	get("a")
	if *calls != before {
		t.Error("a was dropped, want b dropped")
	}
	get("b")
	if *calls != before+1 {
		t.Error("b was kept, want it dropped")
	}
}
