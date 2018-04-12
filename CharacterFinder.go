package arn

// CharacterFinder holds an internal map of ID to anime mappings
// and is therefore very efficient to use when trying to find
// anime by a given service and ID.
type CharacterFinder struct {
	idToCharacter map[string]*Character
}

// NewCharacterFinder creates a new finder for external anime.
func NewCharacterFinder(mappingName string) *CharacterFinder {
	finder := &CharacterFinder{
		idToCharacter: map[string]*Character{},
	}

	for anime := range StreamCharacters() {
		id := anime.GetMapping(mappingName)

		if id != "" {
			finder.idToCharacter[id] = anime
		}
	}

	return finder
}

// GetCharacter tries to find an external anime in our anime database.
func (finder *CharacterFinder) GetCharacter(id string) *Character {
	return finder.idToCharacter[id]
}
