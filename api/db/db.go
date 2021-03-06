package db

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/semyon-dev/hackuniversity/api/model"
	"log"
	"math"
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

	Clicconn, err = sql.Open("clickhouse", "tcp://"+os.Getenv("CLICKHOUSE_HOST")+":9000")
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
			action_time  DateTime,
			to_min Float64 ,
			from_min Float64,
			to_max Float64 ,
			from_max Float64
		) engine=Memory
	`)

	if err != nil {
		log.Println("ошибка при создании таблицы events", err)
	}
}

// создайние event в clickouse
func NewEvent(param string, author string, min, max float64) {

	sqlStr := "SELECT maximum,minimum FROM criticals WHERE paramname =$1"
	rows, err := Conn.Query(sqlStr, param)
	if err != nil {
		fmt.Println(err)
	}

	var oldMin, oldMax float64
	for rows.Next() {
		rows.Scan(&oldMin, &oldMax)
	}

	var (
		tx, _   = Clicconn.Begin()
		stmt, _ = tx.Prepare("INSERT INTO events (PARAM, AUTHOR, action_day, action_time,from_min,to_min,from_max,to_max) VALUES (?, ?, ?, ?,?,?)")
	)
	defer stmt.Close()

	if _, err := stmt.Exec(
		param,
		author,
		time.Now(),
		time.Now(),
		oldMin,
		min,
		oldMax,
		max,
	); err != nil {
		log.Println(err)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
}

func ConnectPostgres() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr = "host=" + os.Getenv("POSTGRES_HOST") + " port=5432 user=semyon dbname=dbtest sslmode=disable"
	//connStr = "host=192.168.1.106 port=5432 user=semyon dbname=dbtest sslmode=disable"

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

	fmt.Println("connected successfully....")
}

func AverageValue(paramName, dateStart, dateEnd string) float64 {
	execStr := "SELECT avg(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "')  AND toDateTime('" + dateEnd + "')"
	rows, err := Clicconn.Query(execStr)
	if err != nil {
		fmt.Println(err)
	}

	var val float64
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			fmt.Println(err)
		}
	}
	return val
}

func MaxValue(paramName, dateStart, dateEnd string) float64 {
	execStr := "SELECT MAX(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "')  AND toDateTime('" + dateEnd + "')"
	rows, err := Clicconn.Query(execStr)
	if err != nil {
		fmt.Println(err)
	}

	var val float64
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			fmt.Println(err)
		}
	}

	return val
}

func MinValue(paramName, dateStart, dateEnd string) float64 {
	execStr := "SELECT MIN(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "')  AND toDateTime('" + dateEnd + "')"
	rows, err := Clicconn.Query(execStr)
	if err != nil {
		fmt.Println(err)
	}

	var val float64
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			fmt.Println(err)
		}
	}

	return val
}

func UpdateCritical(name string, min, max float64) error {
	_, err := Conn.Exec("UPDATE criticals SET minimum = $2,maximum = $3 WHERE paramname = $1", name, min, max)
	return err
}

func GetCriticals() []model.Criticals {
	rows, err := Conn.Query("SELECT paramname,minimum,maximum FROM criticals")
	if err != nil {
		fmt.Println(err)
	}
	var criticals []model.Criticals
	var name string
	var min, max float64
	for rows.Next() {
		err = rows.Scan(&name, &min, &max)
		if err != nil {
			fmt.Println(err)
		}
		criticals = append(criticals, model.Criticals{Name: name, Min: min, Max: max})
	}

	return criticals
}

func GetErrors(dateStart, dateEnd string, limit int) []model.Error {

	var criticalErrors []model.Error

	if dateStart == "today" {
		execStr := "SELECT datetime ,paramName,paramValue,message FROM errors ORDER BY id DESC LIMIT $1"
		rows, err := Conn.Query(execStr, limit)
		if err != nil {
			panic(err)
		}

		var dateTime, paramName, message string
		var paramValue float64
		for rows.Next() {
			err = rows.Scan(&dateTime, &paramName, &paramValue, &message)
			if err != nil {
				fmt.Println(err)
			}
			criticalError := model.Error{Message: message, DateTime: dateTime, ParamValue: paramValue, ParamName: paramName}
			if err != nil {
				fmt.Println(err)
			}
			criticalErrors = append(criticalErrors, criticalError)
			fmt.Println(dateTime, paramName)
		}
	} else {
		execStr := "SELECT datetime,paramName,paramValue,message FROM errors WHERE datetime BETWEEN $1 AND $2 LIMIT $3"
		rows, err := Conn.Query(execStr, dateStart, dateEnd, limit)
		if err != nil {
			panic(err)
		}

		var dateTime, paramName, message string
		var paramValue float64
		for rows.Next() {
			err = rows.Scan(&dateTime, &paramName, &paramValue, &message)
			if err != nil {
				fmt.Println(err)
			}
			criticalError := model.Error{Message: message, DateTime: dateTime, ParamValue: paramValue, ParamName: paramName}
			if err != nil {
				fmt.Println(err)
			}
			criticalErrors = append(criticalErrors, criticalError)
			fmt.Println(dateTime, paramName)
		}
	}

	return criticalErrors
}

func GetHourlyErrors(paramName, date string) []float64 {

	timeStart := model.Time{Hour: 0, Minute: 0, Second: 0}
	timeEnd := model.Time{Hour: 1, Minute: 0, Second: 0}

	var values []float64

	for i := 0; i < 24; i++ {
		dateStart := date + " " + timeStart.ToStringHour()
		dateEnd := date + " " + timeEnd.ToStringHour()

		execStr := "SELECT AVG(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "')  AND toDateTime('" + dateEnd + "')"
		rows, err := Clicconn.Query(execStr)
		if err != nil {
			fmt.Println(err)
		}
		var value float64
		for rows.Next() {
			err := rows.Scan(&value)
			if err != nil {
				fmt.Println(err)
			}
		}
		if math.IsNaN(value) {
			value = 0
		}
		values = append(values, value)
		timeStart.NextHour()
		timeEnd.NextHour()
	}

	//execStr := "SELECT MIN(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
	//rows, err := Clicconn.Query(execStr)
	//if err != nil {
	//	fmt.Println(err)
	//}

	return values
}
