package db

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/semyon-dev/hackuniversity/api/model"
	"log"
	"os"
	"strconv"
	"strings"
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

	fmt.Println("connected successfully....")
}





func AverageValue(paramName, dateStart, dateEnd string) float64 {
	execStr := "SELECT avg(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
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
	execStr := "SELECT MAX(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
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
	execStr := "SELECT MIN(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
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

// unused:
func newDate(date string) model.Date {
	vals := strings.Split(date, ".")

	day, err := strconv.Atoi(vals[0])
	if err != nil {
		fmt.Println(err)
	}
	month, err := strconv.Atoi(vals[0])
	if err != nil {
		fmt.Println(err)
	}
	year, err := strconv.Atoi(vals[0])
	if err != nil {
		fmt.Println(err)
	}
	return model.Date{Day: day, Month: month, Year: year}
}

// unused:
func daysBetween(dateStart, dateEnd model.Date) {
	date1 := time.Date(dateStart.Year, time.Month(dateStart.Month), dateStart.Day, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(dateEnd.Year, time.Month(dateEnd.Month), dateEnd.Day, 0, 0, 0, 0, time.UTC)
	days := int(date2.Sub(date1))
	fmt.Println(days)
}

// unused:
func insertMinMax(name string, min float64, max float64) {
	_, err := Conn.Exec("INSERT INTO criticals(paramname,minimum,maximum) VALUES($1,$2,$3)", name, min, max)
	if err != nil {
		fmt.Println(err)
	}
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


func GetErrors(dateStart,dateEnd string)[]model.Error{
	execStr := "SELECT dateTime,paramName,paramValue,message FROM errors WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
	rows, err := Conn.Query(execStr)
	if err != nil {
		fmt.Println("" + "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		panic(err)
	}

	var criticalErrors []model.Error
	var dateTime,paramName,message string
	var paramValue float32
	for rows.Next() {
		err = rows.Scan(&dateTime,&paramName,&paramValue,&message)
		if err!=nil{
			fmt.Println(err)
		}
		criticalError:=model.Error{Message:message,DateTime:dateTime,ParamValue:paramValue,ParamName:paramName}
		if err != nil {
			fmt.Println(err)
		}
		criticalErrors = append(criticalErrors,criticalError)
	}

	return criticalErrors
}



