package autodocs

import (
	"bufio"
	"os"
	"strings"
)

// Type represents a type in a Go source file.
type Type struct {
	Name    string
	Comment string
}

// Endpoint returns the REST endpoint for that type.
func (typ *Type) Endpoint() string {
	return "/api/" + strings.ToLower(typ.Name) + "/"
}

// GetTypeDocumentation tries to gather documentation about the given type.
func GetTypeDocumentation(typeName string, filePath string) (*Type, error) {
	typ := &Type{
		Name: typeName,
	}

	file, err := os.Open(filePath)

	if err != nil {
		return typ, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var comments []string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		isComment := strings.HasPrefix(line, "// ")

		if isComment {
			comment := strings.TrimPrefix(line, "// ")
			comments = append(comments, comment)
			continue
		}

		if strings.HasPrefix(line, "type ") {
			definedTypeName := strings.TrimPrefix(line, "type ")
			space := strings.Index(definedTypeName, " ")
			definedTypeName = definedTypeName[:space]

			if definedTypeName == typeName {
				typ.Comment = strings.Join(comments, " ")
			}
		}

		if !isComment {
			comments = nil
		}
	}

	if err := scanner.Err(); err != nil {
		return typ, err
	}

	return typ, nil
}
