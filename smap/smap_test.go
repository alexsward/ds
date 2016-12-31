package smap

import (
	"fmt"
	"testing"
)

func TestSmapSize(t *testing.T) {
	fmt.Println("-- TestSmapSize")
	m := New()
	assertSmapSize(t, m, 0)
	if !m.IsEmpty() {
		t.Error("Map size 0 assertion passed, but isEmpty didnt?")
	}
	m.Put("key", "val")
	assertSmapSize(t, m, 1)
}

func TestSmapContains(t *testing.T) {
	fmt.Println("-- TestSmapContains")
	m := New()
	notContains := m.Contains("1")
	if notContains {
		t.Error("Did not expect to contain value that map shouldn't")
	}
	m.Put("1", "a")
	contains := m.Contains("1")
	if !contains {
		t.Error("Expected map to contain value 'a', it didn't")
	}
}

func TestSmapContainsValue(t *testing.T) {
	fmt.Println("-- TestSmapContainsValue")
	m := New()
	m.Put("1", "a")
	contains := m.ContainsValue("a")
	if !contains {
		t.Error("Expected map to contain value 'a', it didn't")
	}
	notContains := m.ContainsValue("whatever")
	if notContains {
		t.Error("Did not expect to contain value that map shouldn't")
	}
}

func TestSmapDelete(t *testing.T) {
	fmt.Println("-- TestSmapDelete")
	m := New()
	m.Put("1", "a")
	assertSmapValue(t, m, "1", "a")
	m.Delete("1")
	val, got := m.Get("1")
	if got {
		t.Errorf("Did not expect a value for key %s, but got %s", "1", val)
	}
}

func TestSmapReplace(t *testing.T) {
	fmt.Println("-- TestSmapReplace")
	m := New()
	replaced := m.Replace("1", "a")
	if replaced {
		t.Error("Did not expect to be able to replace a non-existent key")
	}
	m.Put("1", "a")
	m.Replace("1", "b")
	assertSmapValue(t, m, "1", "b")
}

func TestSmapMerge(t *testing.T) {
	fmt.Println("-- TestSmapMerge")
	m1 := New()
	m1.Put("1", "a")
	m1.Put("2", "b")
	m1.Put("3", "c")
	m2 := New()
	m2.Put("3", "d")
	m2.Put("4", "e")
	m2.Put("5", "f")
	m1.Merge(m2)
	assertSmapValue(t, m1, "1", "a")
	assertSmapValue(t, m1, "2", "b")
	assertSmapValue(t, m1, "3", "d")
	assertSmapValue(t, m1, "4", "e")
	assertSmapValue(t, m1, "5", "f")
}

func TestSmapTransform(t *testing.T) {
	fmt.Println("-- TestSmapTransform")
	m := New()
	m.Put("1", "a")
	m.Put("2", "a")
	m.Put("3", "a")
	m.Transform(func(arg1 string) string {
		return "b"
	})
	assertSmapValue(t, m, "1", "b")
	assertSmapValue(t, m, "2", "b")
	assertSmapValue(t, m, "3", "b")
	i := 1
	m.forEach(func(k, v string) bool {
		if i == 3 {
			return true
		}
		i++
		m.Put(k, "c")
		return false
	})
	assertSmapValue(t, m, "1", "c")
	assertSmapValue(t, m, "2", "c")
	assertSmapValue(t, m, "3", "b")
}

func TestSmapAlter(t *testing.T) {
	fmt.Println("-- TestSmapAlter")
	m := New()
	alter := func(x string) string {
		return "b"
	}
	altered1 := m.Alter("1", alter)
	if altered1 {
		t.Error("Didn't expect to be able to alter key")
	}
	m.Put("1", "a")
	altered2 := m.Alter("1", alter)
	if !altered2 {
		t.Error("Expected to be able to alter key")
	}
}

func assertSmapSize(t *testing.T, m *Map, expected int) {
	if m.Size() != expected {
		t.Errorf("Expected map size to be %d, instead got %d", expected, m.Size())
	}
}

func assertSmapValue(t *testing.T, m *Map, key, value string) {
	got, found := m.Get(key)
	if !found {
		t.Errorf("Expected to find key %s, didn't", key)
	}
	if got != value {
		t.Errorf("Expected key:%s to contain value %s, instead got %s", key, value, got)
	}
}
