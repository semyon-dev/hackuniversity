package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/semyon-dev/hackuniversity/checkerr/log"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var conn *sql.DB

var upgrader = websocket.Upgrader{} // use default options
var Connections []*websocket.Conn

func connect() {
	var err error

	// language=SQL
	connStr := "host=192.168.99.1 port=5432 user=postgres dbname=postgres password=12345678 sslmode=disable"

	conn, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS errors(
 					id SERIAL PRIMARY KEY,
 					dateTime Date,
 					paramName varchar(20),
 					paramValue float8,
 					message text
 					   
		)
`)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

func insertError(name, message string, paramValue float64) error {
	_, err := conn.Exec("INSERT INTO errors(dateTime,paramName,paramValue,message) VALUES(now(),$1,$2,$3)", name, paramValue, message)
	return err
}

type Criticals struct {
	Name string  `json:"param"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

func getCriticals() map[string]map[string]float64 {
	rows, err := conn.Query("SELECT paramname,minimum,maximum FROM criticals")
	if err != nil {
		fmt.Println(err)
	}
	var criticals map[string]map[string]float64
	var name string
	var min, max float64
	for rows.Next() {
		rows.Scan(&name, &min, &max)
		criticals[name] = map[string]float64{"min": min, "max": max}

	}

	return criticals
}

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	Connections = append(Connections, c)
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		checkCriticalParameters(message)

		//for _, Connection := range Connections {
		//	if Connection != c {
		//		err = Connection.WriteMessage(mt, message)
		//		if err != nil {
		//			log.Println("write:", err)
		//			break
		//		}
		//	}
		//}
	}
}

func checkCriticalParameters(jsonData []byte) {
	var data map[string]float64
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err)
	}

	criticals := getCriticals()

	for key, val := range data {
		if (val < criticals[key]["min"]) && (val > criticals[key]["max"]) {
			Log.Error(key + " - over normal value")
			err := insertError(key, "over normal value", val)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	//
	//if data.WATER >= 30 {
	//	Log.Error("параметр WATER превышает норму")
	//	fmt.Println("параметр WATER превышает норму")
	//}
}

func main() {
	connect()
	Logging()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
