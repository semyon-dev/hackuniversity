package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/semyon-dev/hackuniversity/checkerr/db"
	. "github.com/semyon-dev/hackuniversity/checkerr/log"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var addr = flag.String("addr", os.Getenv("LOCAL_IP")+":8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
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
	fmt.Println(os.Getenv("LOCAL_IP")+" -is locip")
	db.Connect()
	Logging()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
