package model

import "strconv"

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

type Time struct {
	Hour int
	Minute int
	Second int
}




func (time *Time)NextHour(){
	time.Hour+=1
}

func (time *Time)ToStringHour()string{
	if time.Hour>9{
		return strconv.Itoa(time.Hour)+":00:00"
	}else {
		return "0"+strconv.Itoa(time.Hour)+":00:00"
	}
}

