package DB

import (
	"database/sql"
	"fmt"
)

var conn *sql.DB


func Connect() {
	var err error

	// language=SQL
	connStr := "host=localhost port=5432 user=postgres dbname=postgres password=12345678 sslmode=disable"

	conn, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS errors(
 					id SERIAL PRIMARY KEY,
 					dateTime Date,
 					paramName varchar(20),
 					paramValue float8,
 					message text
 					   
		)
`)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

func InsertError(name, message string, paramValue float64) error {
	_, err := conn.Exec("INSERT INTO errors(dateTime,paramName,paramValue,message) VALUES(now(),$1,$2,$3)", name, paramValue, message)
	return err
}

type Criticals struct {
	Name string  `json:"param"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

func GetCriticals() map[string]map[string]float64 {
	rows, err := conn.Query("SELECT paramname,minimum,maximum FROM criticals")
	if err != nil {
		fmt.Println(err)
	}
	var criticals map[string]map[string]float64
	var name string
	var min, max float64
	for rows.Next() {
		rows.Scan(&name, &min, &max)
		criticals[name] = map[string]float64{"min": min, "max": max}

	}

	return criticals
}