package websocketAPI

//this is a package for easely setting up a gorilla websocket without needing all the functions yourself. made for quick projects that dont need 100% configuration but still need a simple and configurable websocket
import (
	"fmt"
	"log"
	"net/http"
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//defaultSettings is a type to store all the settings
type settingsBase struct {
	MessageReceived func(c *websocket.Conn, p []byte)                           //handler for received messages can be customized
	MessageReceiver func(c *websocket.Conn, ms settingsBase) //handler for receiving messages can be customized, custom function does not need to use MessageReceived
}

//defaultMessageReceived is the default handler for messages, it just prints them out.
func defaultMessageReceived(c *websocket.Conn, p []byte) {
	fmt.Println(string(p))
}

//defaultMessageReceiver waits for the messages to come and sends them off to settings.MessageReceived
func defaultMessageReceiver(c *websocket.Conn, ms settingsBase) {
	for {
		mt, message, err := c.ReadMessage()

		if err != nil {
			fmt.Println("read:", err)
			break
		}
		ms.MessageReceived(c,message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

//settings stores the ccurrent srttings
var settings = settingsBase{defaultMessageReceived, defaultMessageReceiver}

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
	settings.MessageReceiver(c, settings)
}

var addr = flag.String("addr", "localhost:8080", "http service address")

//TODO: set custom port
//TODO: set custom routes

//Start starts the websocket, all trafic goes through port 8080,
func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", wsUpgrader)
	http.Handle("/", r)
	fmt.Println("[WebsocketAPI] Websocket API is initialized.")
	http.ListenAndServe(*addr, nil)
}

//SetMessageHandler sets the message handler to a custom function that can handle the message.
//this resets the receiver to make sure that the message is handles correctly
func SetMessageHandler(s func(c *websocket.Conn,p []byte)) {
	settings = settingsBase{s, defaultMessageReceiver}
}

//SetReceiverHandler sets the message receiver to a custom function that can handle the connection
func SetReceiverHandler(s func(c *websocket.Conn, ms settingsBase)) {
	settings = settingsBase{settings.MessageReceived, s}
}
