package model

type Date struct {
	Day   int
	Month int
	Year  int
}

type Criticals struct {
	Name   string  `json:"param"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Author string  `json:"author"`
}


type Error struct {
	DateTime string		`json:"dateTime"`
	ParamName string	`json:"paramName"`
	ParamValue float32	`json:"paramValue"`
	Message string		`json:"message"`
}
