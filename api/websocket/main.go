package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"flag"
)

type defaultSettings struct{
	MessageReceived func(p []byte)
	MessageReceiver func(c *websocket.Conn,ms defaultSettings)
}

func defaultMessageReceived(p []byte) {
	fmt.Println(string(p)+"what")
}

func defaultMessageReceiver(c *websocket.Conn,ms defaultSettings) {
	for {
		mt, message, err := c.ReadMessage()

		if err != nil {
			fmt.Println("read:", err)
			break
		}
		ms.MessageReceived(message)
		
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

var Settings = defaultSettings{defaultMessageReceived,defaultMessageReceiver}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsUpgrader(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	Settings.MessageReceiver(c, Settings)
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func Start(){
	r := mux.NewRouter()
	r.HandleFunc("/", wsUpgrader)
	http.Handle("/", r)
	fmt.Println("[WebsocketAPI] Websocket API is initialized.")
	http.ListenAndServe(*addr, nil)
}

func SetMessageHandler(s func(p []byte)){
	Settings = defaultSettings{s,defaultMessageReceiver}
}
func SetReceiverHandler(s func(c *websocket.Conn,ms defaultSettings)){
	Settings = defaultSettings{defaultMessageReceived,s}
}





