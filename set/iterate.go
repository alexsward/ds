package set

import "sync"

var (
	// PredicateAll is a Predicate that matches everything in the set
	PredicateAll = func(item Item) bool {
		return true
	}
)

func (s *set) IterateAll() <-chan Item {
	return s.iterate(PredicateAll, s.Size())
}

func (s *set) Iterate(pred Predicate, limit int) <-chan Item {
	return s.iterate(pred, limit)
}

func (s *set) iterate(pred Predicate, limit int) <-chan Item {
	iter := make(chan Item)

	s.lock.RLock()
	go func(lock *sync.RWMutex) {
		defer lock.RUnlock()
		defer close(iter)
		i := 0
		for _, item := range s.data() {
			if pred(item) {
				iter <- item
				i++
			}
			if i == limit {
				break
			}
		}
	}(s.lock)
	return iter
}
