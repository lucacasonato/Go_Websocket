package customMessageReceiver

import (
	"fmt"
	"github.com/CreativeGuy2013/Go_Websocket/api/websocket"
)

func messageReceived(p []byte) {
	fmt.Println(string(p) + "whut")
}

func main() {
	fmt.Println("running main")
	websocket.SetMessageHandler(messageReceived)
	websocket.Start()
}
