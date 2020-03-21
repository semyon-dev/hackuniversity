package main

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/semyon-dev/hackuniversity/api/model"
	"log"
	"strconv"
	"strings"
	"time"
)

var conn *sql.DB
var clicconn *sql.DB

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	connect()
	clickConnect()

	r.GET("/criticals", func(context *gin.Context) {

		criticals := make(map[string]map[string]float64)

		for _, i := range getCriticals() {
			criticals[i.Name] = map[string]float64{"min": i.Min, "max": i.Max}
		}

		context.JSON(200, criticals)
	})

	r.POST("/critical", func(context *gin.Context) {

		var critical model.Criticals
		err := context.ShouldBindJSON(&critical)
		if err != nil {
			fmt.Println(err)
		}
		err = updateCritical(critical.Name, critical.Min, critical.Max)
		if err != nil {
			context.JSON(500, gin.H{
				"status": "ERROR",
			})
			fmt.Println(err)
		} else {
			context.JSON(200, gin.H{
				"status": "OK",
			})
		}
	})

	// test url: /period?paramName=HUMIDITY&dateStart=2020-03-20&dateEnd=2020-03-30&timeStart=00:00:00&timeEnd=00:00:00
	// return average value between start date and time and end date and time
	r.GET("/average", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		params := averageValue(name, dateTimeStart, dateTimeEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.GET("/max", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		params := maxValue(name, dateTimeStart, dateTimeEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.GET("/min", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		params := minValue(name, dateTimeStart, dateTimeEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.GET("/maindata", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		min := minValue(name, dateTimeStart, dateTimeEnd)
		max := maxValue(name, dateTimeStart, dateTimeEnd)
		avg := averageValue(name, dateTimeStart, dateTimeEnd)

		context.JSON(200,
			gin.H{
				"min": min,
				"avg": avg,
				"max": max,
			})
	})

	err := r.Run(":5000")
	if err != nil {
		fmt.Println("ошибка при запуске API", err)
	}
}

// получение границ даты и времени из юрл
func nameDateTimes(context *gin.Context) (string, string, string) {

	currentTime := time.Now().String()
	strCurrTime := strings.Split(currentTime, ".")[0]
	name := context.Query("paramName")
	dateStart := context.Query("dateStart")
	var dateTimeStart, dateTimeEnd string
	if dateStart == "today" {
		dateTimeStart = strings.Split(strCurrTime, " ")[0] + " 00:00:00"
		dateTimeEnd = strCurrTime
	} else {
		dateEnd := context.Query("dateEnd")
		timeStart := context.Query("timeStart")
		timeEnd := context.Query("timeEnd")

		dateTimeStart = dateStart + " " + timeStart
		dateTimeEnd = dateEnd + " " + timeEnd
	}

	return name, dateTimeStart, dateTimeEnd
}

var connStr = "host=192.168.1.106 port=5432 user=semyon dbname=dbtest sslmode=disable"

func clickConnect() {
	var err error
	clicconn, err = sql.Open("clickhouse", "tcp://192.168.1.109:9000?debug=true")
	if err != nil {
		log.Println("ошибка при подключении к clickhouse", err)
	}
	fmt.Println("-------------------")
	if err := clicconn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println("err", err)
		}
	}
	fmt.Println("-------------------")

	_, err = clicconn.Exec(`
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
}

func connect() {
	var err error
	conn, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.Exec(`
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
	rows, err := conn.Query("SELECT id FROM criticals LIMIT 20")
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
			_, err = conn.Exec(createDefaults, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	fmt.Println("connected successfully....")
}

func averageValue(paramName, dateStart, dateEnd string) float32 {
	execStr := "SELECT avg(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
	fmt.Println(execStr + " - !!!!")
	rows, err := clicconn.Query(execStr)
	if err != nil {
		fmt.Println(err)
	}

	var val float32
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			fmt.Println(err)
		}
	}
	return val
}

func maxValue(paramName, dateStart, dateEnd string) float32 {
	execStr := "SELECT MAX(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
	rows, err := clicconn.Query(execStr)
	if err != nil {
		fmt.Println(err)
	}

	var val float32
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			fmt.Println(err)
		}
	}

	return val
}

func minValue(paramName, dateStart, dateEnd string) float32 {
	execStr := "SELECT MIN(" + paramName + ") FROM journal WHERE action_time BETWEEN toDateTime('" + dateStart + "', 'Europe/Moscow')  AND toDateTime('" + dateEnd + "', 'Europe/Moscow')"
	rows, err := clicconn.Query(execStr)
	if err != nil {
		fmt.Println(err)
	}

	var val float32
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
	_, err := conn.Exec("INSERT INTO criticals(paramname,minimum,maximum) VALUES($1,$2,$3)", name, min, max)
	if err != nil {
		fmt.Println(err)
	}
}

func updateCritical(name string, min, max float64) error {
	_, err := conn.Exec("UPDATE criticals SET minimum = $2,maximum = $3 WHERE paramname = $1", name, min, max)
	return err
}

func getCriticals() []model.Criticals {
	rows, err := conn.Query("SELECT paramname,minimum,maximum FROM criticals")
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
