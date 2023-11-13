package notification

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Socket struct {
	Clients map[string]*websocket.Conn
	Send    func(*websocket.Conn, Event) error
}
type Event struct {
	Name    string
	Message interface{}
}

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all requests
	},
}

func (s *Socket) Broadcast(message Event) {
	for _, client := range s.Clients {
		s.Send(client, message)
	}
}
func (s *Socket) Listen(w http.ResponseWriter, r *http.Request) {
	var id string = r.URL.Query().Get("id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	s.Clients[id] = conn
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				delete(s.Clients, id)
			}
			return
		}
	}
}

// func
func broadcast(conn *websocket.Conn, msg Event) error {
	err := conn.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}
func (s *Socket) CloseHandler(code int, text string) {
}

var ClientSocket Socket = Socket{
	Send:    broadcast,
	Clients: make(map[string]*websocket.Conn),
}
var MerchantSocket Socket = Socket{
	Send:    broadcast,
	Clients: make(map[string]*websocket.Conn),
}
var DriverSocket Socket = Socket{
	Send:    broadcast,
	Clients: make(map[string]*websocket.Conn),
}

func Initsocket() {
	http.HandleFunc("/client", ClientSocket.Listen)
	http.HandleFunc("/merchant", MerchantSocket.Listen)
	http.HandleFunc("/driver", DriverSocket.Listen)
	if os.Getenv("APP_ENV") == "development" {
		log.Fatal(http.ListenAndServe("localhost:8080", nil))
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}
