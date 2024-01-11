package objects

type UserLoginAttempt struct {
	Attempts    int
	LastAttempt int64
}

func (u UserLoginAttempt) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"attempts": u.Attempts,
		"last":     u.LastAttempt,
	}
}

func (u UserLoginAttempt) Increment() IncrementableInterface {
	u.Attempts = u.Attempts + 1
	return u
}

func ReadUserLoginAttemptFromMap(input map[string]interface{}) UserLoginAttempt {
	return UserLoginAttempt{
		Attempts:    input["attempts"].(int),
		LastAttempt: input["last"].(int64),
	}
}

func NewUserLoginAttempt() UserLoginAttempt {
	return UserLoginAttempt{}
}
