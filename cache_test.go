package lrucache

import (
	"fmt"
	"testing"
)

func TestLRUCache_NewCache(t *testing.T) {
	t.Parallel()

	cache := NewCache(5)
	if cache == nil {
		t.Fatal("cache is nil")
	}

	if cache.capacity != 5 {
		t.Fatalf("unexpected capacity. expected: %d | actual: %d", 5, cache.capacity)
	}

	if cache.Len() != 0 {
		t.Fatalf("unexpected map len. expected: %d | actual: %d", 0, cache.Len())
	}

	if cache.head != nil {
		t.Fatal("head is not nil")
	}

	if cache.tail != nil {
		t.Fatal("tail is not nil")
	}
}

func TestLRUCache_Len(t *testing.T) {
	t.Parallel()

	cache := NewCache(5)
	if cache.Len() != 0 {
		t.Fatalf("unexpected map len. expected: %d | actual: %d", 0, cache.Len())
	}

	cache.Set("A", 10)
	cache.Set("B", 20)
	if cache.Len() != 2 {
		t.Fatalf("unexpected map len. expected: %d | got: %d", 2, cache.Len())
	}
}

func TestLRUCache_Set(t *testing.T) {
	t.Parallel()

	cache := NewCache(3)
	cache.Set("A", 10)
	node, ok := cache.items["A"]
	if !ok {
		t.Fatalf("value not found. expected: %d | got: %v", 10, nil)
	}

	if node.Value != 10 {
		t.Fatalf("unexpected value. expected: %d | got %d", 10, node.Value)
	}

	if node.Key != "A" {
		t.Fatalf("unexpected value. expected: %s | got %s", "A", node.Key)
	}

	if cache.head != node {
		t.Fatal("node is not head")
	}

	if cache.tail != node {
		t.Fatal("node is not tail")
	}
}

func TestLRUCache_Get(t *testing.T) {
	t.Parallel()

	cache := NewCache(3)
	val, ok := cache.Get("A")
	if ok {
		t.Fatalf("unexpected found flag. expected: %t | got %t", false, ok)
	}

	if val != -1 {
		t.Fatalf("unexpected value. expected: %d | got %d", -1, val)
	}
}

func TestLRUCache_SetOrder(t *testing.T) {
	t.Parallel()

	cache := NewCache(3)
	cache.Set("A", 10)
	cache.Set("B", 20)
	cache.Set("C", 30)

	if cache.tail.Value != 10 {
		t.Fatalf("unexpected tail value. expected: %d | got: %d", 10, cache.tail.Value)
	}

	if cache.head.Value != 30 {
		t.Fatalf("unexpected head value. expected: %d | got: %d", 30, cache.head.Value)
	}
}

func TestLRUCache_GetUpdatesOrder(t *testing.T) {
	t.Parallel()

	cache := NewCache(3)
	cache.Set("A", 10)
	cache.Set("B", 20)
	cache.Set("C", 30)

	value, ok := cache.Get("A")
	if !ok {
		t.Fatal("A not found")
	}

	if value != 10 {
		t.Fatalf("unexpected value. expected: %d | got: %d", 10, value)
	}

	if cache.head.Key != "A" {
		t.Fatalf("unexpected head. expected: %s | got: %s", "A", cache.head.Key)
	}

	if cache.tail.Key != "B" {
		t.Fatalf("unexpected tail. expected: %s | got: %s", "B", cache.tail.Key)
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	t.Parallel()

	cache := NewCache(3)
	cache.Set("A", 10)
	cache.Set("B", 20)
	cache.Set("C", 30)
	cache.Set("D", 40)

	if cache.Len() != 3 {
		t.Fatalf("unexpected len. expected: %d | got: %d", 3, cache.Len())
	}

	if _, ok := cache.items["A"]; ok {
		t.Fatal("A should be evicted")
	}
}

func TestLRUCache_Concurrency(t *testing.T) {
	t.Parallel()

	cache := NewCache(100)
	done := make(chan struct{}, 100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			cache.Set(fmt.Sprintf("%d", i), i)
			cache.Get(fmt.Sprintf("%d", i))
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}

	if cache.Len() != 100 {
		t.Fatalf(
			"unexpected map len(race condition detected) expected: %d | got: %d",
			100, cache.Len(),
		)
	}
}
