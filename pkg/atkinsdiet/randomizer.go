package atkinsdiet

import (
	"math/rand"
	"sync"
	"time"
)

type randomizer struct {
	rand *rand.Rand
	mu   sync.Mutex
}

func newRandomizer() *randomizer {
	return &randomizer{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Intn returns, as an int8, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func (r *randomizer) Intn(n int) (v int) {
	r.mu.Lock()
	v = r.rand.Intn(n)
	r.mu.Unlock()
	return
}

