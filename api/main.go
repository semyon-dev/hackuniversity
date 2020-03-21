package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var conn *sql.DB

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

	r.Run(":5001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var connStr = "host=localhost port=5432 user=postgres dbname=postgres password=12345678 sslmode=disable"

func connect() {
	var err error
	conn, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	res, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS  criticals (
			id serial primary key, 
			paramname varchar(20),
			maximum float ,
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
	rows, _ := conn.Query("SELECT id FROM criticals LIMIT 1")
	var id int
	for rows.Next() {
		rows.Scan(&id)
		if id == 1 {
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
	fmt.Println(res)
}

func insertMinMax(name string, min float64, max float64) {
	_, err := conn.Exec("INSERT INTO criticals(paramname,minimum,maximum) VALUES($1,$2,$3)", name, min, max)
	if err != nil {
		fmt.Println(err)
	}
}

type Criticals struct {
	Name string  `json:"param"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

func updateCritical(name string, min, max float64) error {

	fmt.Println(name)

	_, err := conn.Exec("UPDATE criticals SET minimum = $2,maximum = $3 WHERE paramname = $1", name, min, max)
	return err
}

func getCriticals() []Criticals {
	rows, err := conn.Query("SELECT paramname,minimum,maximum FROM criticals")
	if err != nil {
		fmt.Println(err)
	}
	var criticals []Criticals
	var name string
	var min, max float64
	for rows.Next() {
		rows.Scan(&name, &min, &max)
		criticals = append(criticals, Criticals{Name: name, Min: min, Max: max})

	}

	fmt.Println(criticals[0].Name + " selected")

	return criticals
}
