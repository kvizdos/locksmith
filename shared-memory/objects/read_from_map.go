package objects

import (
	"fmt"
)

func ReadAllocationFromMap(alloc MemoryAllocation, input map[string]any) (any, error) {
	switch alloc {
	case USER_LOGIN_ATTEMPTS:
		return ReadUserLoginAttemptFromMap(input), nil
	}

	return nil, fmt.Errorf("allocation not supported")
}
