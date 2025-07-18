package common

import (
	"container/heap"
	"testing"
)

func TestAbs(t *testing.T) {
	if Abs(360) != 360 || Abs(-270) != 270 || Abs(int(-90)) != 90 {
		t.Error("Abs failed")
	}
}

func TestPriorityQueue(t *testing.T) {
	type testItem struct {
		value    string
		priority int
	}

	pq := &PriorityQueue[testItem, int]{}
	heap.Init(pq)

	// Push items into the priority queue
	items := []testItem{
		{"item1", 3},
		{"item2", 1},
		{"item3", 2},
		{"item4", 4},
	}
	for _, item := range items {
		heap.Push(pq, &PQItem[testItem, int]{Content: item, Priority: item.priority})
	}

	// Check the length of the priority queue
	if pq.Len() != len(items) {
		t.Errorf("expected length %d, got %d", len(items), pq.Len())
	}

	// Check the top item without removing it
	topItem := pq.GetTop()
	if topItem.Content.value != "item2" {
		t.Errorf("expected top item 'item2', got '%s'", topItem.Content.value)
	}

	// Pop items from the priority queue and check order
	expectedOrder := []string{"item2", "item3", "item1", "item4"}
	for i, expected := range expectedOrder {
		item := heap.Pop(pq).(*PQItem[testItem, int])
		if item.Content.value != expected {
			t.Errorf("expected item '%s' at position %d, got '%s'", expected, i, item.Content.value)
		}
	}

	pq = &PriorityQueue[testItem, int]{}
	heap.Init(pq)

	// Push items into the priority queue
	items = []testItem{
		{"item1", 6},
		{"item2", 8},
	}
	for _, item := range items {
		heap.Push(pq, &PQItem[testItem, int]{Content: item, Priority: item.priority})
	}

	// Check the length of the priority queue
	if pq.Len() != len(items) {
		t.Errorf("expected length %d, got %d", len(items), pq.Len())
	}

	// Check the top item without removing it
	topItem = pq.GetTop()
	if topItem.Content.value != "item1" {
		t.Errorf("expected top item 'item1', got '%s'", topItem.Content.value)
	}

	// Pop items from the priority queue and check order
	expectedOrder = []string{"item1", "item2"}
	for i, expected := range expectedOrder {
		item := heap.Pop(pq).(*PQItem[testItem, int])
		if item.Content.value != expected {
			t.Errorf("expected item '%s' at position %d, got '%s'", expected, i, item.Content.value)
		}
	}
}
