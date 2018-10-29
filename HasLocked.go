package arn

// HasLocked implements common like and unlike methods.
type HasLocked struct {
	Locked bool `json:"locked"`
}

// Lock locks the object.
func (obj *HasLocked) Lock(userID string) {
	obj.Locked = true
}

// Unlock unlocks the object.
func (obj *HasLocked) Unlock(userID string) {
	obj.Locked = false
}

// IsLocked implements the Lockable interface.
func (obj *HasLocked) IsLocked() bool {
	return obj.Locked
}
