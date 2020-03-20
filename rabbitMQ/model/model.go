package model

type Data struct {
	PRESSURE float64 `json:"pressure"`
	HUMIDITY float64 `json:"humidity"`
	TEMP_HOM float64 `json:"temp_hom"`
	TEMP_WOR float64 `json:"temp_wor"`
	LEVEL    float64 `json:"level"`
	MASS     float64 `json:"mass"`
	WATER    float64 `json:"water"`
	CO2      float64 `json:"co_2"`
}
