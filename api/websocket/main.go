package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type MessageReceived interface {
	MessageReceived(p []byte)
}

type settings struct {
	MessageReceived
}

//Settings
var Settings = settings{messageReceived}

func (m MessageReceived) messageReceived(p []byte) {
	fmt.Println(p)
}

//MessageReciver
func (s *settings) MessageReceiver(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()

		if err != nil {
			fmt.Println("read:", err)
			break
		}
		s.messageReceived(message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type rect struct {
	width, height int
}

func (r *rect) area() int {
	return r.width * r.height
}

var upgrader = websocket.Upgrader{} // use default options

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", wsUpgrader)
	http.Handle("/", r)
	fmt.Println("[WebsocketAPI] Websocket API is initialized.")
}

func wsUpgrader(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	Settings.MessageReceiver(c)
}

//HelloWorld is a test function!
func HelloWorld(hello string) {
	fmt.Println(hello)
}

//ParseMessage
func ParseMessage() {

}
