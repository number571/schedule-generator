package schedule

import (
	"time"
	"math/rand"
)

func RandSubgroup() SubgroupType {
	return SubgroupType(random(0, 1))
}

// Returns [min:max] value.
func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min + 1) + min
}
