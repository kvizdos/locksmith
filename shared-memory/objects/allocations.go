package objects

type MemoryAllocation string

const (
	USER_LOGIN_ATTEMPTS MemoryAllocation = "attempts"
)

type MappableInterface interface {
	ToMap() map[string]interface{}
}

type IncrementableInterface interface {
	MappableInterface
	Increment() IncrementableInterface
}

func AllocationFromName(allocation MemoryAllocation) MappableInterface {
	switch allocation {
	case USER_LOGIN_ATTEMPTS:
		return NewUserLoginAttempt()
	default:
		return nil
	}
}
