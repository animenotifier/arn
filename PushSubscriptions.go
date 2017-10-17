package arn

// PushSubscriptions ...
type PushSubscriptions struct {
	UserID UserID              `json:"userId"`
	Items  []*PushSubscription `json:"items"`
}

// Find returns the subscription with the specified ID, if available.
func (list *PushSubscriptions) Find(id string) *PushSubscription {
	for _, item := range list.Items {
		if item.ID() == id {
			return item
		}
	}

	return nil
}

// GetPushSubscriptions ...
func GetPushSubscriptions(id string) (*PushSubscriptions, error) {
	obj, err := DB.Get("PushSubscriptions", id)

	if err != nil {
		return nil, err
	}

	return obj.(*PushSubscriptions), nil
}
