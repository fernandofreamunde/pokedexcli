package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("differenttestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)

			val, ok := cache.Get(c.key)
			if !ok {
				t.Error("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 1
	const waitTime = baseTime + 1
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime * time.Second)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
