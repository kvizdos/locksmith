package providers

import (
	"fmt"
	"sync"

	sharedmemory "github.com/kvizdos/locksmith/shared-memory"
	"github.com/kvizdos/locksmith/shared-memory/objects"
)

type RamSharedMemoryProvider struct {
	// map[Allocation]map[keyName]value
	internalMemory map[objects.MemoryAllocation]map[string]map[string]interface{}
	mutex          sync.RWMutex
}

func NewRamSharedMemoryProvider() *RamSharedMemoryProvider {
	return &RamSharedMemoryProvider{
		internalMemory: make(map[objects.MemoryAllocation]map[string]map[string]interface{}),
	}
}

func (mem *RamSharedMemoryProvider) GetFromMemory(allocation objects.MemoryAllocation, key string) (interface{}, bool) {
	mem.mutex.RLock()
	defer mem.mutex.RUnlock()

	if mem.internalMemory[allocation] == nil {
		mem.internalMemory[allocation] = make(map[string]map[string]interface{})
	}

	v, ok := mem.internalMemory[allocation][key]

	if ok {
		alloc, err := objects.ReadAllocationFromMap(allocation, v)
		if err != nil {
			return nil, false
		}

		return alloc, true
	}

	return nil, false
}

func (mem *RamSharedMemoryProvider) SetMemory(allocation objects.MemoryAllocation, key string, value objects.MappableInterface, options ...sharedmemory.MemoryOptions) error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	if mem.internalMemory[allocation] == nil {
		mem.internalMemory[allocation] = make(map[string]map[string]interface{})
	}

	mem.internalMemory[allocation][key] = value.ToMap()

	return nil
}

func (mem *RamSharedMemoryProvider) DeleteMemory(allocation objects.MemoryAllocation, key string) error {
	delete(mem.internalMemory[allocation], key)
	return nil
}

func (mem *RamSharedMemoryProvider) Increment(allocation objects.MemoryAllocation, key string, alreadyPulled ...interface{}) error {
	var obj interface{}
	found := true

	if len(alreadyPulled) > 0 {
		obj = alreadyPulled[0]
	} else {
		obj, found = mem.GetFromMemory(allocation, key)
	}

	if !found {
		obj = objects.AllocationFromName(allocation)
	}

	if v, ok := obj.(objects.IncrementableInterface); ok {
		err := mem.SetMemory(allocation, key, v.Increment())
		return err
	}

	return fmt.Errorf("allocation does not support incrementing")
}
