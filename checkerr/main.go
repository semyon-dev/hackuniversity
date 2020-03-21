package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/semyon-dev/hackuniversity/checkerr/log"
	"github.com/semyon-dev/hackuniversity/checkerr/model"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

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
	var data model.Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err)
	}
	if data.WATER >= 30 {
		Log.Error("параметр WATER превышает норму")
		fmt.Println("параметр WATER превышает норму")
	}
}

func main() {
	Logging()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}