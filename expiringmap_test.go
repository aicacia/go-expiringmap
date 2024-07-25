package expiringmap

import (
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := New[string, string]()
	if m.Len() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	m := New[string, Animal]()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.Set("elephant", elephant, time.Now().Add(time.Minute))
	m.Set("monkey", monkey, time.Now().Add(time.Minute))

	if m.Len() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsertAbsent(t *testing.T) {
	m := New[string, Animal]()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.SetIfAbsent("elephant", elephant, time.Now().Add(time.Minute))
	if ok := m.SetIfAbsent("elephant", monkey, time.Now().Add(time.Minute)); ok {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {
	m := New[string, Animal]()

	// Get a missing element.
	val, ok := m.Get("Money")

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if (val != Animal{}) {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant, time.Now().Add(time.Minute))

	// Retrieve inserted element.
	elephant, ok = m.Get("elephant")
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestGetOrSet(t *testing.T) {
	m := New[string, Animal]()

	// Set a missing element.
	val := m.GetOrSet("Money", Animal{"elephant"}, time.Now().Add(time.Minute))
	if val.name != "elephant" {
		t.Error("default item was not inserted.")
	}

	// Set a missing element.
	oldVal := m.GetOrSet("Money", Animal{"lion"}, time.Now().Add(time.Minute))
	if oldVal.name != "elephant" {
		t.Error("previous item was not returned")
	}

	if m.Len() != 1 {
		t.Error("map should contain exactly one element.")
	}
}

func TestHas(t *testing.T) {
	m := New[string, Animal]()

	// Get a missing element.
	if m.Has("Money") == true {
		t.Error("element shouldn't exists")
	}

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant, time.Now().Add(time.Minute))

	if m.Has("elephant") == false {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	m := New[string, Animal]()

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey, time.Now().Add(time.Minute))

	m.Remove("monkey")

	if m.Len() != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := m.Get("monkey")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if (temp != Animal{}) {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	m.Remove("noone")
}

func TestCount(t *testing.T) {
	m := New[string, Animal]()
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)}, time.Now().Add(time.Minute))
	}

	if m.Len() != 100 {
		t.Error("Expecting 100 element within map.")
	}
}

func TestAsyncCount(t *testing.T) {
	m := New[string, Animal]()

	var wg sync.WaitGroup
	var size = 100
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(i int) {
			start := i * size
			for j := start; j < start+size; j++ {
				m.Set(strconv.Itoa(j), Animal{strconv.Itoa(j)}, time.Now().Add(time.Minute))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(i int) {
			start := i * size
			for j := start; j < start+size; j++ {
				m.Delete(strconv.Itoa(j))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if m.Len() != 0 {
		t.Errorf("Expecting 0 elements within map, got %d.", m.Len())
	}
}

func TestIsEmpty(t *testing.T) {
	m := New[string, Animal]()

	if m.IsEmpty() == false {
		t.Error("new map should be empty")
	}

	m.Set("elephant", Animal{"elephant"}, time.Now().Add(time.Minute))

	if m.IsEmpty() != false {
		t.Error("map shouldn't be empty.")
	}
}

func TestRange(t *testing.T) {
	m := New[string, Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)}, time.Now().Add(time.Minute))
	}

	counter := 0
	// Iterate over elements.
	m.Range(func(key string, value Animal) bool {
		counter++
		return true
	})

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestIterator(t *testing.T) {
	m := New[string, Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)}, time.Now().Add(time.Minute))
	}

	counter := 0
	// Iterate over elements.
	for item := range m.Iter() {
		val := item.Val

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestBufferedIterator(t *testing.T) {
	m := New[string, Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)}, time.Now().Add(time.Minute))
	}

	counter := 0
	// Iterate over elements.
	for item := range m.Iter() {
		val := item.Val

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestClear(t *testing.T) {
	m := New[string, Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)}, time.Now().Add(time.Minute))
	}

	m.Clear()

	if m.Len() != 0 {
		t.Error("We should have 0 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	m := New[int, int]()
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Set(i, i, time.Now().Add(time.Minute))

			// Retrieve item from map.
			val, _ := m.Get(i)

			// Write to channel inserted value.
			ch <- val
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Set(i, i, time.Now().Add(time.Minute))

			// Retrieve item from map.
			val, _ := m.Get(i)

			// Write to channel inserted value.
			ch <- val
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])

	// Make sure map contains 1000 elements.
	if m.Len() != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestKeys(t *testing.T) {
	m := New[int, Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(i, Animal{strconv.Itoa(i)}, time.Now().Add(time.Minute))
	}

	count := 0
	for range m.Keys() {
		count += 1
	}
	if count != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestExpire(t *testing.T) {
	m := New[string, Animal]()

	m.Set("elephant", Animal{"elephant"}, time.Now().Add(-time.Minute))

	if m.IsEmpty() != true {
		t.Error("map should be empty.")
	}
}
