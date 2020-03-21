// by Semyon

package main

import (
	"github.com/semyon-dev/hackuniversity/pusher/db"
	"github.com/semyon-dev/hackuniversity/pusher/websocket"
)

func main() {
	db.Connect()
	websocket.Run()
}
