package connect

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"log"
)

var conn *sql.DB

func Connect() {
	var err error
	conn, err = sql.Open("clickhouse", "tcp://192.168.1.106:8123?debug=true")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			id serializable primary key, 
			PRESSURE  float,
			HUMIDITY  float,
			TEMP_HOME  float,
			TEMP_WORK  float,
			LEVELS  float,
			MASS  float,
			WATER  float,
			CO2 float 
		) 
	`)
	//engine=Memory
	if err != nil {
		panic(err)
	}
	fmt.Println("connected successfully....")
}

func SaveAll(data map[string]float32) {

	var (
		tx, _   = conn.Begin()
		// rewrite for standart
		stmt, _ = tx.Prepare("INSERT INTO example (t1,t2,t3,t4,t5,t6,t7) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()

	stmt.Exec(
		data["PRESSURE"],
		data["HUMIDITY"],
		data["TEMP_HOME"],
		data["TEMP_WORK"],
		data["LEVELS"],
		data["MASS"],
		data["WATER"],
		data["CO2"],
	)

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
