package cmsketch

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func createTestCMSketch() CMSketch {
	cm, _ := New(0.999, 0.001)
	return cm
}

func TestCMSketchCreate(t *testing.T) {
	fmt.Println("TestCMSketchCreate")

	cases := []struct {
		delta        float64
		epsilon      float64
		depth, width uint
	}{
		{0.9, 0.01, 4, 272},
		{0.99, 0.001, 7, 2719},
		{0.999, 0.0001, 10, 27183},
		{0.9999, 0.0001, 14, 27183},
		{0.99999, 0.0001, 17, 27183},
	}

	for _, c := range cases {
		cm, _ := New(c.delta, c.epsilon)
		assert(t, uint64(c.width), uint64(cm.Width()), fmt.Sprintf("epsilon:%f delta:%f expected width:%d actual:%d", c.epsilon, c.delta, c.width, cm.Width()))
		assert(t, uint64(c.depth), uint64(cm.Depth()), fmt.Sprintf("epsilon:%f delta:%f expected depth:%d actual:%d", c.epsilon, c.delta, c.depth, cm.Depth()))
	}
}

func TestCMSketchAddSingle(t *testing.T) {
	fmt.Println("TestCMSketchAddSingle")

	cm := createTestCMSketch()
	cm.Add([]byte("Alex"), 1)
	c := cm.Count([]byte("Alex"))
	assert(t, 1, c, "Count mismatched")
}

func TestCMSketchAddMultiple(t *testing.T) {
	fmt.Println("TestCMSketchAddMultiple")

	cm := createTestCMSketch()
	cm.Add([]byte("Alex"), 10)
	c := cm.Count([]byte("Alex"))
	assert(t, 10, c, "Count mismatched")
}

func TestParallelAdds(t *testing.T) {
	fmt.Println("TestParallelAdds")

	cm := createTestCMSketch()

	var wg sync.WaitGroup
	adder := func(t uint64) {
		cm.Add([]byte("Alex"), t)
		wg.Done()
	}

	wg.Add(5)
	go adder(1)
	go adder(5)
	go adder(10)
	go adder(3)
	go adder(9)
	wg.Wait()

	c := cm.Count([]byte("Alex"))
	assert(t, 28, c, "Count mismatched")
}

func TestCMSketchRemove(t *testing.T) {
}

func TestCMSketchCount(t *testing.T) {
}

func TestCMSketchHash(t *testing.T) {
}

func TestCMSketchMerge(t *testing.T) {
	fmt.Println("TestCMSketchMerge")

	cm1 := createTestCMSketch()
	cm2 := createTestCMSketch()

	cm1.Add([]byte("Alex"), 1)
	cm2.Add([]byte("Alex"), 2)

	c1 := cm1.Count([]byte("Alex"))
	assert(t, 1, c1, "Count mismatched")

	c2 := cm2.Count([]byte("Alex"))
	assert(t, 2, c2, "Count mismatched")

	cm1.Merge(cm2)

	c := cm1.Count([]byte("Alex"))
	assert(t, 3, c, "Count mismatched")
}

func TestCMSketchMergeUneven(t *testing.T) {
	fmt.Println("TestCMSketchMergeUneven")

	cm1 := createTestCMSketch()
	cm2, _ := New(0.99, 0.0001)
	err := cm1.Merge(cm2)
	if err != ErrCannotMergeDifferentDimensions {
		t.Errorf("Expected error %s, got [%s] instead", ErrCannotMergeDifferentDimensions, err)
	}
}

func TestSize(t *testing.T) {
	fmt.Println("TestSize")

	cm := createTestCMSketch()
	if cm.Size() != cm.Depth()*cm.Width() {
		t.Error("Expected size to be equal to depth x width")
	}
}

func assert(t *testing.T, expected, actual uint64, msg string) {
	if actual != expected {
		t.Errorf("%s -- expected:[%d] actual:[%d]", msg, expected, actual)
	}
}

func benchmark(b *testing.B, runes []byte, length int) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	generator := func() byte {
		i := r.Intn(len(runes))
		return runes[i]
	}

	cm := createTestCMSketch()
	for i := 0; i < b.N; i++ {
		var key []byte
		for j := 0; j < length; j++ {
			key = append(key, generator())
		}

		cm.Add(key, 1)
	}
}

func BenchmarkAddFew6(b *testing.B) {
	fmt.Println("BenchmarkAddFew6")

	runes := []byte("abcdefgh123456789")
	benchmark(b, runes, 6)
}

func BenchmarkAddFew12(b *testing.B) {
	fmt.Println("BenchmarkAddFew6")

	runes := []byte("abcdefgh123456789")
	benchmark(b, runes, 12)
}

func BenchmarkAddMany6(b *testing.B) {
	fmt.Println("BenchmarkAddMany")

	runes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	benchmark(b, runes, 6)
}

func BenchmarkAddMany12(b *testing.B) {
	runes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	benchmark(b, runes, 12)
}

func BenchmarkAddMany25(b *testing.B) {
	runes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	benchmark(b, runes, 12)
}

func BenchmarkAddMany45(b *testing.B) {
	runes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	benchmark(b, runes, 12)
}
