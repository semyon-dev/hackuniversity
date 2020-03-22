package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/semyon-dev/hackuniversity/checkerr/db"
	. "github.com/semyon-dev/hackuniversity/checkerr/log"
	"github.com/semyon-dev/hackuniversity/checkerr/model"
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
	var data model.Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err)
	}

	criticals := db.GetCriticals()
	checkValues("TEMPWORK", criticals, data.TEMPWORK)
	checkValues("TEMPHOME", criticals, data.TEMPHOME)
	checkValues("WATER", criticals, data.WATER)
	checkValues("PRESSURE", criticals, data.PRESSURE)
	checkValues("MASS", criticals, data.MASS)
	checkValues("LEVELPH", criticals, data.LEVELPH)
	checkValues("LEVELCO2", criticals, data.LEVELCO2)
	checkValues("HUMIDITY", criticals, data.HUMIDITY)
}

func checkValues(name string, dict map[string]map[string]float64, value float64) {
	if (value < dict[name]["min"]) || (value > dict[name]["max"]) {
		Log.Error(name + " - превышено допустимое значение")
		err := db.InsertError(name, name+" - превышено допустимое значение", value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	db.Connect()
	Logging()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println("произошла ошибка при запуске вебсокетов:", err)
	}
}
