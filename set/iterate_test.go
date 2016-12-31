package set

import (
	"fmt"
	"strconv"
	"testing"
)

var testPredicate = func(i Item) bool {
	num, err := strconv.Atoi(i.Key())
	if err != nil {
		panic(err)
	}
	return num%2 == 1
}

// TestIterateAll verifies you can iterate an entire set
func TestIterateAll(t *testing.T) {
	fmt.Println("TestIterateAll")
	s := getPopulatedSet(1, 5)
	var actual []string
	for i := range s.IterateAll() {
		actual = append(actual, i.Key())
	}
	assertKeyArrayEquals(t, []string{"1", "2", "3", "4", "5"}, actual)
}

// TestIterateLimit tests iterating a limited number of items
func TestIterateLimit(t *testing.T) {
	fmt.Println("TestIterateLimit")
	s := getPopulatedSet(1, 10)
	var actual []string
	for i := range s.Iterate(PredicateAll, 5) {
		actual = append(actual, i.Key())
	}
	if len(actual) != 5 {
		t.Errorf("Expected the length of returned items %d, got %d", 5, len(actual))
	}
}

// TestIteratePredicate tests iterating with a predicate
func TestIteratePredicate(t *testing.T) {
	fmt.Println("TestIteratePredicate")
	s := getPopulatedSet(1, 10)
	testIterate(t, s, testPredicate, s.Size(), []string{"1", "3", "5", "7", "9"})
}

// TestIteratePredicateLimit tests iterating with a predicate and limit
func TestIteratePredicateLimit(t *testing.T) {
	fmt.Println("TestIteratePredicateLimit")
	s := getPopulatedSet(1, 20)
	var actual []string
	for i := range s.Iterate(testPredicate, 5) {
		actual = append(actual, i.Key())
	}
	if len(actual) != 5 {
		t.Errorf("Expected length of returned results to be 5, got %d", len(actual))
	}
	for i, item := range actual {
		if !testPredicate(s.data()[item]) {
			t.Errorf("Returned item [%s] at index %d didn't pass predicate", item, i)
		}
	}
}

func testIterate(t *testing.T, s Set, p Predicate, limit int, expected []string) {
	var actual []string
	for i := range s.Iterate(p, limit) {
		actual = append(actual, i.Key())
	}
	assertKeyArrayEquals(t, expected, actual)
}
