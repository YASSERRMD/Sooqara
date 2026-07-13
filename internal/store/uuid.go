package store

import (
	"crypto/rand"
	"fmt"
)

// newUUID generates a random hex ID.
func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("rand read: %v", err))
	}
	return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
