package arn

// Company ...
type Company struct {
	ID          string      `json:"id"`
	Name        CompanyName `json:"name"`
	Image       string      `json:"image"`
	Description string      `json:"description"`
	Location    Location    `json:"location"`
	Mappings    []*Mapping  `json:"mappings"`
	Created     string      `json:"created"`
	CreatedBy   string      `json:"createdBy"`
	Edited      string      `json:"edited"`
	EditedBy    string      `json:"editedBy"`
}
