package set

import "sync"

// Set is a data structure that will contain all unqiue Items
type Set interface {
	Get(Item) (Item, bool)
	Add(Item) bool
	Remove(Item) bool
	Contains(Item) bool
	Empty()
	Size() int
	IterateAll() <-chan Item
	Iterate(pred Predicate, limit int) <-chan Item

	data() map[string]Item
}

// Item is something in the Set
type Item interface {
	Key() string
}

// Comparison is the result of a Comparator
type Comparison int

const (
	lessThan    = Comparison(-1)
	equal       = Comparison(0)
	greaterThan = Comparison(1)
)

// Comparator compares two items for equality
type Comparator func(i1, i2 Item) Comparison

// NewSet returns a Set that om
func NewSet() Set {
	s := &set{
		m:    make(map[string]Item),
		lock: &sync.RWMutex{},
	}
	return s
}

type set struct {
	lock *sync.RWMutex
	m    map[string]Item
}

func (s *set) data() map[string]Item {
	return s.m
}

func (s *set) Get(item Item) (Item, bool) {
	if value, exists := s.m[item.Key()]; exists {
		return value, exists
	}
	return nil, false
}

func (s *set) Add(item Item) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.m[item.Key()]; exists {
		return false
	}

	s.m[item.Key()] = item
	return true
}

func (s *set) Remove(item Item) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.m[item.Key()]; !exists {
		return false
	}

	delete(s.m, item.Key())
	return true
}

func (s *set) Contains(item Item) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, exists := s.m[item.Key()]
	return exists
}

func (s *set) Empty() {
	s.m = make(map[string]Item)
}

func (s *set) Size() int {
	return len(s.m)
}

// Predicate is used for evaluating items for various operations
type Predicate func(Item) bool

// Union returns the union of s1 and s2 such that the values of keys of s1 are taken first
func Union(s1, s2 Set) (Set, error) {
	union := NewSet()
	for _, item := range s1.data() {
		union.Add(item)
	}
	for _, item := range s2.data() {
		if union.Contains(item) {
			continue
		}
		union.Add(item)
	}

	return union, nil
}

// Intersection returns a set containing mutual items of s1 and s2
func Intersection(s1, s2 Set) (Set, error) {
	intersection := NewSet()

	larger, smaller := s1, s2
	if s2.Size() > s1.Size() {
		larger = s2
		smaller = s1
	}
	for key, item := range larger.data() {
		if _, exists := smaller.data()[key]; !exists {
			continue
		}
		intersection.Add(item)
	}
	return intersection, nil
}

// Difference s1-s2 is all elements in s1 that aren't in s2
func Difference(s1, s2 Set) (Set, error) {
	diff := NewSet()
	for _, item := range s1.data() {
		if !s2.Contains(item) {
			diff.Add(item)
		}
	}

	return diff, nil
}

// Filter returns a set with only the items matching the predicate from all the sets
func Filter(p Predicate, many ...Set) (Set, error) {
	filtered := NewSet()
	for _, s := range many {
		for _, item := range s.data() {
			if filtered.Contains(item) {
				continue
			}
			if p(item) {
				filtered.Add(item)
			}
		}
	}
	return filtered, nil
}

// Equal returns true if two sets have the same Items
func Equal(s1, s2 Set) bool {
	// TODO: locks?
	if len(s1.data()) != len(s2.data()) {
		return false
	}

	for _, item := range s1.data() {
		if !s2.Contains(item) {
			return false
		}
	}

	return true
}

// Subset determines if s2 is a subset of s1
func Subset(s1, s2 Set) bool {
	if s2.Size() > s1.Size() {
		return false
	}

	if s2.Size() == s1.Size() {
		return Equal(s1, s2)
	}

	for _, item := range s2.data() {
		if !s1.Contains(item) {
			return false
		}
	}

	return true
}

// Superset determines if s2 is a super set of s1
func Superset(s1, s2 Set) bool {
	if s2.Size() < s1.Size() {
		return false
	}

	if s2.Size() == s1.Size() {
		return Equal(s1, s2)
	}

	for _, item := range s1.data() {
		if !s2.Contains(item) {
			return false
		}
	}

	return true
}

// DeepEqual is like Equal but also makes sure the sets have the same values
func DeepEqual(s1, s2 Set, comp Comparator) bool {
	if s1.Size() != s2.Size() {
		return false
	}

	for item := range s1.IterateAll() {
		item2, exists := s2.Get(item)
		if !exists || comp(item, item2) != equal {
			return false
		}
	}

	return true
}
