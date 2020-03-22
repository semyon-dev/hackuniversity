package db

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var Conn *sql.DB
var Clicconn *sql.DB
var connStr string

func ConnectClickhouse() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr = "host=" + os.Getenv("POSTGRES_HOST") + " port=5432 user=semyon dbname=dbtest sslmode=disable"

	Clicconn, err = sql.Open("clickhouse", "tcp://"+os.Getenv("CLICKHOUSE_HOST")+":9000?debug=true")
	if err != nil {
		log.Println("ошибка при подключении к clickhouse", err)
	}
	fmt.Println("-------------------")
	if err := Clicconn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println("err", err)
		}
	}
	fmt.Println("-------------------")

	// журнал всех событий от opc server
	_, err = Clicconn.Exec(`
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
		log.Println("ошибка при создании таблицы journal", err)
	}

	// журнал событий изменений критических параметров
	_, err = Clicconn.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			PARAM  String,
			AUTHOR String,
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)

	if err != nil {
		log.Println("ошибка при создании таблицы events", err)
	}
}

// создайние event в clickouse
func NewEvent(param string, author string) {
	var (
		tx, _   = Clicconn.Begin()
		stmt, _ = tx.Prepare("INSERT INTO events (PARAM, AUTHOR, action_day, action_time) VALUES (?, ?, ?, ?)")
	)
	defer stmt.Close()

	if _, err := stmt.Exec(
		param,
		author,
		time.Now(),
		time.Now(),
	); err != nil {
		log.Println(err)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
}

func ConnectPostgres() {
	var err error
	Conn, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	_, err = Conn.Exec(`
		CREATE TABLE IF NOT EXISTS criticals (
			id serial primary key, 
			paramname varchar(20),
			maximum float,
			minimum float 
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	listNames := []string{
		"PRESSURE",
		"HUMIDITY",
		"TEMPHOME",
		"TEMPWORK",
		"LEVELPH",
		"MASS",
		"WATER",
		"LEVELCO2",
	}

	var haveDefaults = false
	rows, err := Conn.Query("SELECT id FROM criticals LIMIT 20")
	if err != nil {
		fmt.Println(err)
	}
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
		}
		if id != 1 {
			haveDefaults = true
		}
	}

	if !haveDefaults {
		for _, i := range listNames {
			createDefaults := "INSERT INTO criticals(paramname,minimum,maximum) VALUES ($1,2,98)"
			_, err = Conn.Exec(createDefaults, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	fmt.Println("connected successfully....")
}
