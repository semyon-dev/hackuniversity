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

// old not used version
// deprecated:
func checkCriticalParameters2(jsonData []byte) {
	var data model.Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err)
	}

	criticals := db.GetCriticals()

	if (data.HUMIDITY < criticals["HUMIDITY"]["min"]) || (data.HUMIDITY > criticals["HUMIDITY"]["max"]) {
		Log.Error("HUMIDITY - over normal value")
		err := db.InsertError("HUMIDITY", "over normal value", data.HUMIDITY)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.LEVELCO2 < criticals["LEVELCO2"]["min"]) || (data.LEVELCO2 > criticals["LEVELCO2"]["max"]) {
		Log.Error("LEVELCO2 - over normal value")
		err := db.InsertError("LEVELCO2", "over normal value", data.LEVELCO2)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.LEVELPH < criticals["LEVELPH"]["min"]) || (data.LEVELPH > criticals["LEVELPH"]["max"]) {
		Log.Error("LEVELPH - over normal value")
		err := db.InsertError("LEVELPH", "over normal value", data.LEVELPH)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.MASS < criticals["MASS"]["min"]) || (data.MASS > criticals["MASS"]["max"]) {
		Log.Error("MASS - over normal value")
		err := db.InsertError("MASS", "over normal value", data.MASS)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.PRESSURE < criticals["PRESSURE"]["min"]) || (data.PRESSURE > criticals["PRESSURE"]["max"]) {
		Log.Error("PRESSURE - over normal value")
		err := db.InsertError("PRESSURE", "over normal value", data.PRESSURE)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.WATER < criticals["WATER"]["min"]) || (data.WATER > criticals["WATER"]["max"]) {
		Log.Error("WATER - over normal value")
		err := db.InsertError("WATER", "over normal value", data.WATER)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.TEMPHOME < criticals["TEMPHOME"]["min"]) || (data.TEMPHOME > criticals["TEMPHOME"]["max"]) {
		Log.Error("TEMPHOME - over normal value")
		err := db.InsertError("TEMPHOME", "over normal value", data.TEMPHOME)
		if err != nil {
			fmt.Println(err)
		}
	}

	if (data.TEMPWORK < criticals["TEMPWORK"]["min"]) || (data.TEMPWORK > criticals["TEMPWORK"]["max"]) {
		Log.Error("TEMPWORK - over normal value")
		err := db.InsertError("TEMPWORK", "over normal value", data.TEMPWORK)
		if err != nil {
			fmt.Println(err)
		}
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
		Log.Error(" - over normal value")
		err := db.InsertError(name, "over normal value", value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	fmt.Println(os.Getenv("LOCAL_IP"))
	db.Connect()
	Logging()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
