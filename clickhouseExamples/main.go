package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/semyon-dev/hackuniversity/clickhouse/connect"
	"github.com/semyon-dev/hackuniversity/clickhouse/websock"
)

func main() {

	Connect()

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
