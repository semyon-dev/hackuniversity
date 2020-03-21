package main
import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	connect()

	r.POST("/criticals", func(context *gin.Context) {

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}



var conn *sql.DB

func connect() {
	var err error
	conn, err = sql.Open("clickhouse", "tcp://192.168.1.106:8123?debug=true")
	if err != nil {
		panic(err)
	}

	res, err := conn.Exec(`
		CREATE TABLE  criticals (
			id serializable primary key, 
			paramname varchar(20),
			maximum float ,
			minimum float ,
			
		) 
	`)

	//engine=Memory
	if err != nil {
		panic(err)
	}

	fmt.Println("connected successfully....")
	fmt.Println(res)
}


func insertMinMax(name string,min float64,max float64){
	_,err:=conn.Exec("INSERT INTO criticals(paramname,minimum,maximum) VALUES($1,$2,$3)",name,min,max)
	if err!=nil{
		fmt.Println(err)
	}
}

