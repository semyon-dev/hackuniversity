package websock

import (
	"awesomeProject4/DB"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "localhost:5000", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func wsConnect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade sended:", err)
		return
	}
	defer c.Close()

	go sendData(c)

	//for {
	//	mt, message, err := c.ReadMessage()
	//	if err != nil {
	//		log.Println("read:", err)
	//		break
	//	}
	//	log.Printf("recv: %s", message)
	//	err = c.WriteMessage(mt, message)
	//	if err != nil {
	//		log.Println("write:", err)
	//		break
	//	}
	//}
}

// send data recuirsively
func sendData(conn *websocket.Conn){

	err:= conn.WriteJSON(DB.GetLastData())
	if err!=nil{
		fmt.Println("err websock")
		fmt.Println(err)
		return
	}else {
		time.Sleep(3000)
		sendData(conn)
	}
}







func TestWS() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/getData",wsConnect)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

