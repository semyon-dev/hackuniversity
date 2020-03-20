package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)


type Skills struct {
	It   []string `json:"it"`
	Home string   `json:"home"`
}

type Jusers struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Sskills Skills `json:"skills"`
}

func Test() {

	app := gin.Default()

	connStr := "host=localhost port=5432 user=postgres dbname=postgres password=12345678 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		fmt.Println("-----------")
	}

	// Language=SQL
	//clickhouse.Exec("drop table users")

	var initdb = "create table if not exists users ( id SERIAL PRIMARY KEY,us_name VARCHAR(30), age INT )"
	_, err = db.Exec(initdb)
	if err != nil {
		fmt.Println(err)
	}

	//
	//_,err = clickhouse.Exec("INSERT INTO users(us_name, age) VALUES ('QWERTY',23)")
	//if err!=nil{
	//	fmt.Println( string(err.Error()) +"  qwertyukgfjh")
	//}

	// Language=SQL
	test3 := "select skills -> 'skills' -> 'it' -> 0 from jusers where cast( skills->'skills'->>'home' as varchar) ='chilling'"
	query3, err := db.Query(test3)
	if err != nil {
		panic(err)
	}

	var tstr string
	for query3.Next() {
		query3.Scan(&tstr)
		println(tstr + " - choto")
	}

	// test clickhouse 2
	// Language=SQL
	var initdb2 = "CREATE TABLE IF NOT EXISTS jusers(id SERIAL PRIMARY KEY, name VARCHAR(30), skills json)"
	_, err = db.Exec(initdb2)
	if err != nil {
		fmt.Println(err)
	}
	//it's shit

	var tstruct = Jusers{Age: 18, Name: "jake", Sskills: Skills{Home: "chilling", It: []string{".net", "go"}}}

	var cnvrted, _ = json.Marshal(tstruct)

	//check sqlx
	// insert raw in $1
	_, err = db.Exec("INSERT INTO jusers(name,skills) VALUES ('semen',$1)", string(cnvrted))
	if err != nil {
		panic(err)
	}

	// Language=SQL
	contains := "SELECT skills->>'hex'::text FROM jusers "
	rows, err := db.Query(contains)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(" rows ---------")
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			panic(err)
		}
		fmt.Println(value)
	}
	fmt.Println(" end rows ---------")

	app.GET("/", func(context *gin.Context) {
		fmt.Println("some --------")
		var dict = make(map[string]string)
		dict["message"] = "ok"
		rows, _ := db.Query("select * from users")
		fmt.Println(rows)

		context.JSON(200, dict)
	})

	err = app.Run(":8090")

	if err != nil {
		log.Panic(err)
	}
}

func showRows(rows *sql.Rows) {
	var value string
	fmt.Println(" rows ---------")
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			panic(err)
		}
		fmt.Println(value)
	}
	fmt.Println(" end rows ---------")
}

