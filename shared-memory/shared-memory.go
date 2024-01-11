package sharedmemory

import (
	"time"

	"github.com/kvizdos/locksmith/shared-memory/objects"
)

type MemoryOptions struct {
	ValidUntil time.Duration // not used yet
}

type MemoryProvider interface {
	GetFromMemory(allocation objects.MemoryAllocation, key string) (interface{}, bool)
	SetMemory(allocation objects.MemoryAllocation, key string, value objects.MappableInterface, opts ...MemoryOptions) error
	DeleteMemory(allocation objects.MemoryAllocation, key string) error
	Increment(allocation objects.MemoryAllocation, key string, alreadyPulled ...interface{}) error
}
