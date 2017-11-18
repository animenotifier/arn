package arn

// Location ...
type Location struct {
	CountryName string  `json:"countryName"`
	CountryCode string  `json:"countryCode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CityName    string  `json:"cityName"`
	RegionName  string  `json:"regionName"`
	TimeZone    string  `json:"timeZone"`
	ZipCode     string  `json:"zipCode"`
}

// IPInfoDBLocation ...
type IPInfoDBLocation struct {
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	CityName    string `json:"cityName"`
	RegionName  string `json:"regionName"`
	TimeZone    string `json:"timeZone"`
	ZipCode     string `json:"zipCode"`
}
