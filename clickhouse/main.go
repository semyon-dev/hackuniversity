package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/hackuniversity/clickhouse/websock"
)

func main() {

	//var (
	//	tx, _   = connect.Begin()
	//	stmt, _ = tx.Prepare("INSERT INTO example (country_code, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
	//)
	//defer stmt.Close()
	//
	//for i := 0; i < 100; i++ {
	//	if _, err := stmt.Exec(
	//		"RU",
	//		10+i,
	//		100+i,
	//		clickhouse.Array([]int16{1, 2, 3}),
	//		time.Now(),
	//		time.Now(),
	//	); err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//
	//if err := tx.Commit(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//rows, err := connect.Query("SELECT country_code, os_id, browser_id, categories, action_day, action_time FROM example")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var (
	//		country               string
	//		os, browser           uint8
	//		categories            []int16
	//		actionDay, actionTime time.Time
	//	)
	//	if err := rows.Scan(&country, &os, &browser, &categories, &actionDay, &actionTime); err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("country: %s, os: %clickhouse, browser: %clickhouse, categories: %v, action_day: %s, action_time: %s", country, os, browser, categories, actionDay, actionTime)
	//}
	//
	//if err := rows.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if _, err := connect.Exec("DROP TABLE example"); err != nil {
	//	log.Fatal(err)
	//}

	go Connect()

	websock.TestWS()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

}
