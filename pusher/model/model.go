package model

type Data struct {
	PRESSURE float64 `json:"PRESSURE"`
	HUMIDITY float64 `json:"HUMIDITY"`
	TEMPHOME float64 `json:"TEMPHOME"`
	TEMPWORK float64 `json:"TEMPWORK"`
	LEVELPH  float64 `json:"LEVELPH"`
	MASS     float64 `json:"MASS"`
	WATER    float64 `json:"WATER"`
	LEVELCO2 float64 `json:"LEVELCO2"`
	MESSAGE  string  `json:"MESSAGE"`
}
