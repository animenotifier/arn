package arn

// IDList ...
type IDList []string

// GetIDList ...
func GetIDList(id string) (IDList, error) {
	obj, err := DB.Get("IDList", id)

	if err != nil {
		return nil, err
	}

	return *obj.(*IDList), nil
}

// Append appends the given ID to the end of the list.
func (idList IDList) Append(id string) {
	idList = append(idList, id)
}
