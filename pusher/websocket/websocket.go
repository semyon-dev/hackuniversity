package websocket

import (
	"flag"
	"github.com/semyon-dev/hackuniversity/pusher/db"
	"github.com/semyon-dev/hackuniversity/pusher/opc"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", os.Getenv("LOCAL_IP")+":8080", "http service address")
var addrDimaWS = flag.String("addrDimaWS", "192.168.1.109:8080", "http service address")

func Run() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	uDima := url.URL{Scheme: "ws", Host: *addrDimaWS, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("fatal error: dial:", err)
	}
	defer c.Close()

	cDima, _, err := websocket.DefaultDialer.Dial(uDima.String(), nil)
	if err != nil {
		log.Println("error: dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			dataJson, data := opc.GetData()
			// отправялемь данные сразу и в бд и в вебсокет
			db.Save(data)
			err := c.WriteMessage(websocket.TextMessage, dataJson)
			if err != nil {
				log.Println("write:", err)
				return
			}
			// отправляем вебскоеты на другой микросервис (Диме)
			err = cDima.WriteMessage(websocket.TextMessage, dataJson)
			_ = t
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
