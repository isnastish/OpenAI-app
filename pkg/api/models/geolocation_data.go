package models

type GeolocationData struct {
	Ip          string `json:"ip"`
	City        string `json:"city"`
	Region      string `json:"region"`
	RegionCode  string `json:"region_code"`
	Country     string `json:"country_name"`
	CountryCode string `json:"country_code"`
	// error handling
	ErrorCode string `json:"code,omitempty"`
	ErrorMsg  string `json:"error,omitempty"`
}
