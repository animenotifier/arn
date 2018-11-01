package arn

// HasText includes a text field.
type HasText struct {
	Text string `json:"text" editable:"true" type:"textarea"`
}

// GetText returns the text of the object.
func (obj *HasText) GetText() string {
	return obj.Text
}
