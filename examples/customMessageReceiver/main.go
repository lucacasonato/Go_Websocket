package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/CreativeGuy2013/Go_Websocket/api/websocket"
)

func messageReceived(c *websocket.Conn,p []byte) {
	fmt.Println(string(p) + "whut")
}

func main() {
	fmt.Println("running main")
	websocketAPI.SetMessageHandler(messageReceived)
	websocketAPI.Start()
}
