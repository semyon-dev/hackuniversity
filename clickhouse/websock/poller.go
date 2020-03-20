package websock

import (
	"database/sql"
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

var conn *sql.DB

func GetLastData() map[string]float32 {
	var data map[string]float32
	rows, err := conn.Query("SELECT t1,t2,t3,t4,t5,t6,t7 FROM example Where id=(select max(id) from example)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		t1, t2, t3, t4, t5, t6, t7 float32
	)

	for rows.Next() {

		if err := rows.Scan(&t1, &t2, &t3, &t4, &t5, &t6, &t7); err != nil {
			log.Fatal(err)
		}
		log.Printf("t1: %f, t2: %f, t3: %f, t4: %f, t5: %f, t6: %f,t7:%f", t1, t2, t3, t4, t5, t6, t7)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	data["t1"] = t1
	data["t2"] = t2
	data["t3"] = t3
	data["t4"] = t4
	data["t5"] = t5
	data["t6"] = t6
	data["t7"] = t7

	return data
}

// send data recuirsively
func sendData(conn *websocket.Conn) {

	err := conn.WriteJSON(GetLastData())
	if err != nil {
		fmt.Println("err websock")
		fmt.Println(err)
		return
	} else {
		time.Sleep(3000)
		sendData(conn)
	}
}

func TestWS() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/getData", wsConnect)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
