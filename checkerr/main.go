package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/semyon-dev/hackuniversity/checkerr/db"
	. "github.com/semyon-dev/hackuniversity/checkerr/log"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var conn *sql.DB

var upgrader = websocket.Upgrader{} // use default options
var Connections []*websocket.Conn



func echo(w http.ResponseWriter, r *http.Request) {
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
	}
}

func checkCriticalParameters(jsonData []byte) {
	var data map[string]float64
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err)
	}

	criticals := db.GetCriticals()

	for key, val := range data {
		if (val < criticals[key]["min"]) && (val > criticals[key]["max"]) {
			Log.Error(key + " - over normal value")
			err := db.InsertError(key, "over normal value", val)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	db.Connect()

	Logging()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
