package models

type Geolocation struct {
	Ip          string `json:"ip,omitempty"`
	City        string `json:"city,omitempty"`
	Region      string `json:"region,omitempty"`
	RegionCode  string `json:"region_code,omitempty"`
	Country     string `json:"country_name,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	// error handling
	ErrorCode string `json:"code,omitempty"`
	ErrorMsg  string `json:"error,omitempty"`
}
