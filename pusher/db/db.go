package db

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/semyon-dev/hackuniversity/pusher/model"
	"log"
	"time"
)

var conn *sql.DB

func Connect() {
	var err error
	conn, err = sql.Open("clickhouse", "tcp://192.168.1.106:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-------------------")
	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println("err", err)
		}
	}
	fmt.Println("-------------------")

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS journal (
			PRESSURE   Float64,
			HUMIDITY Float64,
			TEMPHOME Float64,
			TEMPWORK Float64,
			LEVELPH Float64,
			MASS Float64,
			WATER Float64,
			LEVELCO2 Float64,
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func Save(data model.Data) {
	var (
		tx, _   = conn.Begin()
		stmt, _ = tx.Prepare("INSERT INTO journal (PRESSURE, HUMIDITY, TEMPHOME,TEMPWORK, LEVELPH , MASS, WATER, LEVELCO2, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()

	if _, err := stmt.Exec(
		data.PRESSURE,
		data.HUMIDITY,
		data.TEMPHOME,
		data.TEMPWORK,
		data.LEVELPH,
		data.MASS,
		data.WATER,
		data.LEVELCO2,
		time.Now(),
		time.Now(),
	); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
