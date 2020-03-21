package main

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

var conn *sql.DB
var clicconn *sql.DB

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	connect()

	r.GET("/criticals", func(context *gin.Context) {

		criticals := make(map[string]map[string]float64)

		for _, i := range getCriticals() {
			criticals[i.Name] = map[string]float64{"min": i.Min, "max": i.Max}
		}

		context.JSON(200, criticals)
	})

	r.POST("/critical", func(context *gin.Context) {

		var critical Criticals
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

	r.GET("/period", func(context *gin.Context) {
		name := context.Query("paramName")
		dateStart := context.Query("dateStart")
		dateEnd := context.Query("dateEnd")
		params := getParamForPeriod(name, dateStart, dateEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.Run(":5001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var connStr = "host=192.168.1.106 port=5432 user=semyon dbname=dbtest sslmode=disable"

func clickConnect() {
	var err error
	clicconn, err = sql.Open("clickhouse", "tcp://192.168.1.109:9000?debug=true")
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

func connect() {
	var err error
	conn, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS criticals (
			id serial primary key, 
			paramname varchar(20),
			maximum float,
			minimum float 
		)
	`)

	//engine=Memory
	if err != nil {
		panic(err)
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
	rows, _ := conn.Query("SELECT id FROM criticals LIMIT 20")
	var id int
	for rows.Next() {
		rows.Scan(&id)
		if id != 1 {
			haveDefaults = true
		}
	}

	if !haveDefaults {
		for _, i := range listNames {
			createDefaults := "INSERT INTO criticals(paramname,minimum,maximum) VALUES ($1,2,98)"
			conn.Exec(createDefaults, i)
		}
	}

	fmt.Println("connected successfully....")
}
