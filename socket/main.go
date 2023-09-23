package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mohamedabdifitah/processor/service"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Socket struct {
	Conn *websocket.Conn
	Id   string
}

// socket BaseSockeServer
type Server struct {
	Users  map[string]*Socket
	Events map[string]func(socket *Socket, message []byte)
}

func (s *Server) listen(w http.ResponseWriter, r *http.Request) {
	var id string = r.URL.Query().Get("id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrade this connection to a WebSocket
	// connection
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	var socket *Socket = &Socket{
		Conn: Conn,
		Id:   id,
	}
	s.Users[id] = socket
	fmt.Println(s.Users)
	s.reader(socket)
}
func (s *Server) On(event string, handler func(socket *Socket, message []byte)) {
	s.Events[event] = handler
}
func (s Server) reader(socket *Socket) {
	for {
		_, p, err := socket.Conn.ReadMessage()
		if err != nil {
			continue
		}
		var v struct {
			Event   string      `json:"event"`
			Message interface{} `json:"message"`
		}
		err = json.Unmarshal(p, &v)
		if err != nil {
			log.Println(err)
		}
		// search for Events
		handler, ok := s.Events[v.Event]
		fmt.Println("message: ", v.Message)
		if !ok {
			// println("finding handler of event")
			if err := socket.Conn.WriteJSON(v.Message); err != nil {
				log.Println(err)
			}
			return
		}
		byt, err := json.Marshal(v.Message)
		if err != nil {
			fmt.Println(err)
		}
		handler(socket, byt)
	}

}
func NewServer() *Server {
	return &Server{
		Users:  make(map[string]*Socket),
		Events: make(map[string]func(socket *Socket, message []byte)),
	}
}

var BaseSockeServer *Server

func InitSocket() {
	BaseSockeServer = NewServer()
	BaseSockeServer.On("location", func(socket *Socket, message []byte) {
		var location struct {
			Name        string    `json:"name"`
			Coordinates []float64 `json:"coordinates"`
		}
		err := json.Unmarshal(message, &location)
		if err != nil {
			fmt.Print(err)
		}
		res, err := service.SetDriverLocation(fmt.Sprintf(location.Name+":"+socket.Id), location.Coordinates[0], location.Coordinates[1])
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println(res)
	})
	http.HandleFunc("/ws", BaseSockeServer.listen)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
