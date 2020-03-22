package model

type Date struct {
	Day   int
	Month int
	Year  int
}

type Criticals struct {
	Name string  `json:"param"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}
