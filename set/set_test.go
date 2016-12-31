package set

import (
	"fmt"
	"strconv"
	"testing"
)

// TODO: more appropriate to test the implementations, not the interface

// TestGet
func TestGet(t *testing.T) {
	fmt.Println("TestGet")
	s := getTestSet()
	one := testItem{1}
	s.Add(testItem{1})
	if item, found := s.Get(one); !found {
		t.Errorf("Test item not found, expected it to be")
	} else if item.Key() != one.Key() {
		t.Errorf("Expected returned item %s, got back %s", one.Key(), item.Key())
	}

	if _, found := s.Get(testItem{2}); found {
		t.Errorf("Did not expect to get an item back")
	}
}

// TestAdd validates adding items to a set
func TestAdd(t *testing.T) {
	fmt.Println("TestAdd")
	s := getTestSet()
	s.Add(testItem{1})
}

// TestRemove validates set removal
func TestRemove(t *testing.T) {
	fmt.Println("TestRemove")
	s := getTestSet()
	assertOperation(t, "add 1 to set", s.Add(testItem{1}), true)
	assertSetSize(t, s, 1)
	assertOperation(t, "remove 1 from set", s.Remove(testItem{1}), true)
	assertSetSize(t, s, 0)
	assertOperation(t, "remove 1 from set", s.Remove(testItem{1}), false)
	assertSetSize(t, s, 0)
}

// TestContains verifies the contains functionality works
func TestContains(t *testing.T) {
	fmt.Println("TestContains")
	s := getTestSet()
	i := testItem{1}
	assertOperation(t, "empty set contains 1", s.Contains(i), false)
	s.Add(i)
	assertOperation(t, "empty set contains 1", s.Contains(i), true)
	s.Remove(i)
	assertOperation(t, "empty set contains 1", s.Contains(i), false)
}

// TestEmpty validates emptying a set
func TestEmpty(t *testing.T) {
	fmt.Println("TestEmpty")
	s := getTestSet()
	assertOperation(t, "new set is empty", s.Size() == 0, true)
	populateSet(s, 1, 10)
	assertOperation(t, "new set is not empty after population empty", s.Size() == 10, true)
	s.Empty()
	assertOperation(t, "new set is empty after Empty() operation", s.Size() == 0, true)
}

// TestSize validates set sizing
func TestSize(t *testing.T) {
	fmt.Println("TestSize")
	s := getTestSet()
	for i := 1; i < 26; i++ {
		s.Add(testItem{i})
	}
	assertSetSize(t, s, 25)
	for i := 1; i < 26; i++ {
		s.Add(testItem{i})
	}
	assertSetSize(t, s, 25)
	for i := 1; i < 11; i++ {
		s.Remove(testItem{i})
	}
	assertSetSize(t, s, 15)
}

// TestUnion1
func TestUnion1(t *testing.T) {
	fmt.Println("TestUnion1")
	s1, s2 := getTestSet(), getTestSet()
	populateSet(s1, 1, 10)
	populateSet(s2, 11, 20)
	assertSetContains(t, s1, 1, 10)
	assertSetContains(t, s2, 11, 20)
	u, err := Union(s1, s2)
	if err != nil {
		t.Error(err)
	}
	assertSetContains(t, u, 1, 20)
}

func TestUnion2(t *testing.T) {
	fmt.Println("TestUnion2")
	s1, s2 := getTestSet(), getTestSet()
	populateSet(s1, 1, 10)
	populateSet(s2, 1, 15)
	u, err := Union(s1, s2)
	if err != nil {
		t.Error(err)
	}
	assertSetSize(t, u, 15)
	assertSetContains(t, u, 1, 15)
}

func TestUnion3(t *testing.T) {
	fmt.Println("TestUnion2")
	f := func(delta int) func(int) int {
		return func(i int) int {
			return i + delta
		}
	}
	s1, s2 := getTestSet(), getTestSet()
	populateSetItem2(s1, 1, 10, f(100))
	populateSetItem2(s2, 1, 10, f(200))
	u, err := Union(s1, s2)
	if err != nil {
		t.Error(err)
	}
	assertSetSize(t, u, 10)
	assertSetContains(t, u, 1, 10)
	for i := 1; i < 11; i++ {
		item, found := u.Get(testItem{i})
		if !found {
			t.Errorf("Item with key %s not found in set", item.Key())
		}
		transformed := f(100)(i)
		if item.(testItem2).extra != transformed {
			t.Errorf("Expected Item with key to have transformed extra %d, instead got %d", transformed, item.(testItem2).extra)
		}
	}
}

// TestIntersection
func TestIntersection(t *testing.T) {
	fmt.Println("TestIntersection")
	var tests = []struct {
		size                           int
		s1start, s1end, s2start, s2end int
		containsStart, containsEnd     int
	}{
		{5, 1, 5, 1, 5, 1, 5},
		{3, 1, 5, 1, 3, 1, 3},
		{3, 1, 3, 1, 5, 1, 3},
		{0, 1, 5, 6, 10, 1, 0},
	}
	for i, test := range tests {
		s1, s2 := getPopulatedSet(test.s1start, test.s1end), getPopulatedSet(test.s2start, test.s2end)
		x, err := Intersection(s1, s2)
		if err != nil {
			t.Errorf("Error with case %d: %s", i, err)
		}
		assertSetSize(t, x, test.size)
		assertSetContains(t, x, test.containsStart, test.containsEnd)
	}
}

func TestEqual1(t *testing.T) {
	fmt.Println("TestEqual1")
	s1, s2 := getTestSet(), getTestSet()
	populateSet(s1, 1, 10)
	populateSet(s2, 1, 10)
	if !Equal(s1, s2) {
		t.Errorf("Sets should be equal")
	}
}

func TestEqual2(t *testing.T) {
	fmt.Println("TestEqual2")
	s1, s2 := getTestSet(), getTestSet()
	populateSet(s1, 1, 5)
	populateSet(s2, 1, 10)
	if Equal(s1, s2) {
		t.Errorf("Sets should not be equal, they are different lengths")
	}
}

func TestEqual3(t *testing.T) {
	fmt.Println("TestEqual3")
	s1, s2 := getTestSet(), getTestSet()
	populateSet(s1, 1, 10)
	populateSet(s2, 11, 20)
	if Equal(s1, s2) {
		t.Errorf("Sets should not be equal, they have different items")
	}
}

func TestDeepEqual(t *testing.T) {
	fmt.Println("TestDeepEqual")
	var tests = []struct {
		equal  bool
		s1, s2 Set
	}{
		{false, getPopulatedSet(1, 2), getPopulatedSet(1, 3)},
		{false, getPopulatedSet(1, 2), getPopulatedSet(3, 4)},
		{false, getPopulatedSet(1, 2), getPopulatedSet(1, 4)},
		{true, getPopulatedSet(1, 2), getPopulatedSet(1, 2)},
		{true, getPopulatedSet(1, 0), getPopulatedSet(1, 0)},
	}
	comp := func(i1, i2 Item) Comparison {
		n1, err := strconv.Atoi(i1.Key())
		if err != nil {
			panic(err)
		}
		n2, err := strconv.Atoi(i2.Key())
		if err != nil {
			panic(err)
		}
		if n1 != n2 {
			return lessThan // doesn't conform to the interface, but we only care about equality for this test
		}
		return equal
	}
	for i, test := range tests {
		if DeepEqual(test.s1, test.s2, comp) != test.equal {
			t.Errorf("Case %d failed, expected equality of %t, got %t", i+1, test.equal, !test.equal)
		}
	}
}

func TestSubset(t *testing.T) {
	fmt.Println("TestSubset")
	var tests = []struct {
		subset bool
		s1, s2 Set
	}{
		{true, getPopulatedSet(1, 10), getPopulatedSet(1, 10)},
		{true, getPopulatedSet(1, 10), getPopulatedSet(1, 5)},
		{true, getPopulatedSet(1, 10), getPopulatedSet(1, 0)},
		{false, getPopulatedSet(1, 10), getPopulatedSet(1, 20)},
		{false, getPopulatedSet(1, 10), getPopulatedSet(11, 15)},
	}

	for i, test := range tests {
		subset := Subset(test.s1, test.s2)
		if subset != test.subset {
			t.Errorf("Case %d failed, expected subset %t, got %t", i+1, test.subset, subset)
		}
	}
}

func TestSuperset(t *testing.T) {
	fmt.Println("TestSuperset")
	var tests = []struct {
		superset bool
		s1, s2   Set
	}{
		{true, getPopulatedSet(1, 10), getPopulatedSet(1, 10)},
		{true, getPopulatedSet(1, 10), getPopulatedSet(1, 20)},
		{true, getPopulatedSet(5, 15), getPopulatedSet(1, 20)},
		{false, getPopulatedSet(5, 10), getPopulatedSet(1, 5)},
		{false, getPopulatedSet(5, 10), getPopulatedSet(1, 9)},
	}
	for i, test := range tests {
		sup := Superset(test.s1, test.s2)
		if sup != test.superset {
			t.Errorf("Case %d failed, expected superset %t, got %t", i, test.superset, sup)
		}
	}
}

func TestDifference(t *testing.T) {
	fmt.Println("TestDifference")
	var tests = []struct {
		setArgs, expected []int
	}{
		{setArgs: []int{1, 10, 1, 10}, expected: []int{}},
		{setArgs: []int{1, 10, 1, 5}, expected: []int{6, 7, 8, 9, 10}},
		{setArgs: []int{1, 5, 1, 10}, expected: []int{}},
	}
	for i, test := range tests {
		s1 := getPopulatedSet(test.setArgs[0], test.setArgs[1])
		s2 := getPopulatedSet(test.setArgs[2], test.setArgs[3])
		d, err := Difference(s1, s2)
		if err != nil {
			t.Errorf("Error with case %d: %s", i, err)
		}
		assertSetSize(t, d, len(test.expected))
		assertSetContainsItems(t, d, test.expected)
	}
}

func TestFilter(t *testing.T) {
	fmt.Println("TestFilter")
	var tests = []struct {
		expected []string
		many     []Set
		p        Predicate
	}{
		{[]string{"1", "2", "3", "4", "5"}, []Set{getPopulatedSet(1, 5)}, func(i Item) bool { return true }},
		{[]string{}, []Set{getPopulatedSet(1, 5)}, func(i Item) bool { return false }},
		{[]string{"1", "5", "3", "7", "9"}, []Set{getPopulatedSet(1, 5), getPopulatedSet(1, 10)}, testPredicate},
	}
	for i, test := range tests {
		f, err := Filter(test.p, test.many...)
		if err != nil {
			t.Errorf("Error filtering for case %d: %s", i, err)
		}
		var actual []string
		for x := range f.IterateAll() {
			actual = append(actual, x.Key())
		}
		assertKeyArrayEquals(t, actual, test.expected)
	}
}

type testItem struct {
	i int
}

type testItem2 struct {
	testItem
	extra int
}

func getTestSet() Set {
	return NewSet()
}

func getPopulatedSet(b, e int) Set {
	s := getTestSet()
	populateSet(s, b, e)
	return s
}

func populateSet(s Set, start, end int) {
	for i := start; i <= end; i++ {
		s.Add(testItem{i})
	}
}

func populateSetItem2(s Set, start, end int, f func(int) int) {
	for i := start; i <= end; i++ {
		s.Add(testItem2{testItem{i}, f(i)})
	}
}

func (ti testItem) Key() string {
	return strconv.Itoa(ti.i)
}

func assertSetSize(t *testing.T, s Set, expected int) {
	if s.Size() != expected {
		t.Errorf("Expected size = %d, instead got %d", expected, s.Size())
	}
}

func assertOperation(t *testing.T, op string, result, expected bool) {
	if result != expected {
		t.Errorf("operation: %s failed, expected %t, got %t", op, expected, result)
	}
}

func assertSetContains(t *testing.T, s Set, start, end int) {
	for i := start; i <= end; i++ {
		if !s.Contains(testItem{i}) {
			t.Errorf("set should contain item %d, it did not", i)
		}
	}
}

func assertSetContainsItems(t *testing.T, s Set, expected []int) {
	for _, expectation := range expected {
		i := testItem{expectation}
		if !s.Contains(i) {
			t.Errorf("Expected set to contain item %s, it did not", i.Key())
		}
	}
}

func assertKeyArrayEquals(t *testing.T, expected, actual []string) {
	if len(expected) != len(actual) {
		t.Errorf("Lengths of expected and actual arrays not equal: %d vs %d", len(expected), len(actual))
	}
	for i, expect := range expected {
		found := false
		for j := range actual {
			if actual[j] == expect {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Case %d failed, expected to find %s and didn't", i+1, expect)
		}
	}
}
