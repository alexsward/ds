package cmsketch

import (
	"encoding/binary"
	"errors"
	"hash/fnv"
	"math"
	"sync"
)

// CMSketch is the interface for interacting with an underlying Count-Min Sketch
// Create this interface using the New function exposed by this package
type CMSketch interface {
	Add(item []byte, total uint64) error
	Remove(item []byte, total uint64) error
	Count(item []byte) uint64
	Merge(merge CMSketch) error
	Width() uint
	Depth() uint
	Size() uint
	data(x, y int) uint64
}

// implementation struct for CMSketch
type cmsketch struct {
	d, w  uint
	grid  [][]uint64
	h     hasher
	mutex sync.RWMutex
}

type hasher func(item []byte) []uint

var (
	// ErrCannotMergeDifferentDimensions when CMSketch.Merge is invoked with another CMSketch with different dimensions
	ErrCannotMergeDifferentDimensions = errors.New("cannot merge CMSketch objects with different dimensions")
)

// New returns a new CountMin Sketch
func New(delta, epsilon float64) (CMSketch, error) {
	cm := &cmsketch{
		d: depth(delta),
		w: width(epsilon),
	}

	cm.grid = make([][]uint64, cm.d)
	for i := 0; uint(i) < cm.d; i++ {
		cm.grid[i] = make([]uint64, cm.w)
	}

	cm.h = newHasher(cm.d, cm.w)

	return cm, nil
}

func newHasher(d, w uint) hasher {
	h := fnv.New64()
	return func(item []byte) []uint {
		defer h.Reset()
		h.Write(item)

		u, l := upperAndLower(h.Sum(nil))
		positions := make([]uint, d)
		for i := uint(0); i < d; i++ {
			position := (u*i + l) % w
			positions[i] = position
		}
		return positions
	}
}

func upperAndLower(hashed []byte) (uint, uint) {
	u := binary.BigEndian.Uint32(hashed[0:4])
	l := binary.BigEndian.Uint32(hashed[4:8])
	return uint(u), uint(l)
}

// calculates the width of the sketch's grid
func width(episilon float64) uint {
	r := math.E / episilon
	return uint(math.Ceil(r))
}

// calculates the depth of the sketch's grid
func depth(delta float64) uint {
	r := math.Log(1-delta) / math.Log(0.5)
	return uint(math.Ceil(r))
}

// Add places an item into the data structure
func (cms *cmsketch) Add(item []byte, count uint64) error {
	cms.lock(true)
	defer cms.unlock(true)

	ls := cms.getLocations(item)
	for i := uint(0); i < cms.Depth(); i++ {
		y := ls[i]
		cms.grid[i][y] += count
	}

	return nil
}

// Remove drops the count of an item, if possible
func (cms *cmsketch) Remove(item []byte, count uint64) error {
	cms.lock(true)
	defer cms.unlock(true)

	ls := cms.getLocations(item)
	for i := uint(0); i < cms.Depth(); i++ {
		y := ls[i]
		cms.grid[i][y] -= count
	}

	return nil
}

// Count returns the count of the item
func (cms *cmsketch) Count(item []byte) uint64 {
	cms.lock(false)
	defer cms.unlock(false)

	var min uint64
	for x, y := range cms.getLocations(item) {
		if min == 0 || cms.grid[x][y] < min {
			min = cms.grid[x][y]
		}
	}

	return min
}

// Merge combines two CMSketch stuctures, if they're equivalently sized
// If they're not equivalently sized, returns ErrCannotMergeDifferentDimensions
func (cms *cmsketch) Merge(merge CMSketch) error {
	cms.lock(true)
	defer cms.unlock(true)

	if cms.Depth() != merge.Depth() || cms.Width() != merge.Width() {
		return ErrCannotMergeDifferentDimensions
	}

	for r, row := range cms.grid {
		for c := range row {
			cms.grid[r][c] += merge.data(r, c)
		}
	}

	return nil
}

func (cms *cmsketch) Width() uint {
	return cms.w
}

func (cms *cmsketch) Depth() uint {
	return cms.d
}

func (cms *cmsketch) Size() uint {
	return cms.d * cms.w
}

func (cms *cmsketch) data(x, y int) uint64 {
	return cms.grid[x][y]
}

func (cms *cmsketch) lock(write bool) {
	if write {
		cms.mutex.Lock()
		return
	}

	cms.mutex.RLock()
}

func (cms *cmsketch) unlock(write bool) {
	if write {
		cms.mutex.Unlock()
		return
	}

	cms.mutex.RUnlock()
}

func (cms *cmsketch) getLocations(item []byte) []uint {
	return cms.h(item)
}
