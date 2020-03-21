package model

type Data struct {
	PRESSURE float64 `json:"pressure"`
	HUMIDITY float64 `json:"humidity"`
	TEMPHOME float64 `json:"temp_hom"`
	TEMPWORK float64 `json:"temp_wor"`
	LEVELPH    float64 `json:"level"`
	MASS     float64 `json:"mass"`
	WATER    float64 `json:"water"`
	LEVELCO2 float64 `json:"levelco_2"`
}