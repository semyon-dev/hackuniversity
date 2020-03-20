package DB

import (
	"database/sql"
	"fmt"
	"log"

	_ 	"github.com/ClickHouse/clickhouse-go"

)

var conn *sql.DB


func Connect(){
	var err error
	conn, err = sql.Open("clickhouse", "tcp://192.168.1.106:8123?debug=true")
	if err!=nil{
		panic(err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			t1  float,
			t2  float,
			t3  float,
			t4  float,
			t5  float,
			t6  float,
			t7  float
		) 
	`)
	//engine=Memory
	if err!=nil{
		panic(err)
	}
	fmt.Println("connected successfully....")
}


func saveAll(data map[string]float32){

	var (
		tx, _   = conn.Begin()
		stmt, _ = tx.Prepare("INSERT INTO example (t1,t2,t3,t4,t5,t6,t7) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()


	stmt.Exec(
		data["t1"],
		data["t2"],
		data["t3"],
		data["t4"],
		data["t5"],
		data["t6"],
		data["t7"],
		);

	if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}


}
