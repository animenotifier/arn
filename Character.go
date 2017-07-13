package arn

// Character ...
type Character struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

// GetCharacter ...
func GetCharacter(id string) (*Character, error) {
	obj, err := DB.Get("Character", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Character), nil
}

// StreamCharacters returns a stream of all characters.
func StreamCharacters() (chan *Character, error) {
	objects, err := DB.All("Character")
	return objects.(chan *Character), err
}

// MustStreamCharacters returns a stream of all characters.
func MustStreamCharacters() chan *Character {
	stream, err := StreamCharacters()
	PanicOnError(err)
	return stream
}

// AllCharacters returns a slice of all characters.
func AllCharacters() ([]*Character, error) {
	var all []*Character

	stream, err := StreamCharacters()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
