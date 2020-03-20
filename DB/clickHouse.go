package DB

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
	if err != nil {
		panic(err)
	}
	fmt.Println("connected successfully....")
}

func saveAll(data map[string]float32) {

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
	)

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}

func GetLastData() map[string]float32 {
	var data map[string]float32
	rows, err := conn.Query("SELECT t1,t2,t3,t4,t5,t6,t7 FROM example Where id=(select max(id) from example)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		t1, t2, t3, t4, t5, t6, t7 float32
	)

	for rows.Next() {

		if err := rows.Scan(&t1, &t2, &t3, &t4, &t5, &t6, &t7); err != nil {
			log.Fatal(err)
		}
		log.Printf("t1: %f, t2: %f, t3: %f, t4: %f, t5: %f, t6: %f,t7:%f", t1, t2, t3, t4, t5, t6, t7)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	data["t1"] = t1
	data["t2"] = t2
	data["t3"] = t3
	data["t4"] = t4
	data["t5"] = t5
	data["t6"] = t6
	data["t7"] = t7

	return data
}
