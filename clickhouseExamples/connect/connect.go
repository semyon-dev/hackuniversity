package connect

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"log"
	"time"
)

var conn *sql.DB

func Connect() {
	connect, err := sql.Open("clickhouse", "tcp://192.168.1.106:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--------------")
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println("err", err)
		}
	}
	fmt.Println("--------------")

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			country_code FixedString(2),
			os_id        UInt8,
			browser_id   UInt8,
			categories   Array(Int16),
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)

	if err != nil {
		log.Fatal(err)
	}
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO example (country_code, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		if _, err := stmt.Exec(
			"RU",
			10+i,
			100+i,
			clickhouse.Array([]int16{1, 2, 3}),
			time.Now(),
			time.Now(),
		); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	rows, err := connect.Query("SELECT country_code, os_id, browser_id, categories, action_day, action_time FROM example")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			country               string
			os, browser           uint8
			categories            []int16
			actionDay, actionTime time.Time
		)
		if err := rows.Scan(&country, &os, &browser, &categories, &actionDay, &actionTime); err != nil {
			log.Fatal(err)
		}
		log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_day: %s, action_time: %s", country, os, browser, categories, actionDay, actionTime)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	//
	//if _, err := connect.Exec("DROP TABLE example"); err != nil {
	//	log.Fatal(err)
	//}
}

func SaveAll(data map[string]float32) {

	var (
		tx, _ = conn.Begin()
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
