package arn

// ActivityCreate is a user activity that creates something.
type ActivityCreate struct {
	ObjectType string `json:"objectType"`
	ObjectID   string `json:"objectId"`

	HasID
	HasCreator
}

// NewActivityCreate creates a new activity.
func NewActivityCreate(objectType string, objectID string, userID string) *ActivityCreate {
	return &ActivityCreate{
		HasID: HasID{
			ID: GenerateID("ActivityCreate"),
		},
		HasCreator: HasCreator{
			Created:   DateTimeUTC(),
			CreatedBy: userID,
		},
		ObjectType: objectType,
		ObjectID:   objectID,
	}
}

// Object returns the object that was created.
func (activity *ActivityCreate) Object() interface{} {
	obj, _ := DB.Get(activity.ObjectType, activity.ObjectID)
	return obj
}

// Postable casts the object to the Postable interface.
func (activity *ActivityCreate) Postable() Postable {
	return activity.Object().(Postable)
}

// Type returns the type name.
func (activity *ActivityCreate) Type() string {
	return "ActivityCreate"
}
