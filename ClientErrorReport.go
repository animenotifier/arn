package arn

// ClientErrorReport saves JavaScript errors that happen in web clients like browsers.
type ClientErrorReport struct {
	ID           string `json:"id"`
	Message      string `json:"message"`
	Stack        string `json:"stack"`
	FileName     string `json:"fileName"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`

	HasCreator
}
