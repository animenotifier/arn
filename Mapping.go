package arn

// Mapping ...
type Mapping struct {
	Service   string `json:"service"`
	ServiceID string `json:"serviceId"`
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}

// Name ...
func (mapping *Mapping) Name() string {
	switch mapping.Service {
	case "shoboi/anime":
		return "Shoboi"
	default:
		return ""
	}
}

// Link ...
func (mapping *Mapping) Link() string {
	switch mapping.Service {
	case "shoboi/anime":
		return "http://cal.syoboi.jp/tid/" + mapping.ServiceID
	default:
		return ""
	}
}
