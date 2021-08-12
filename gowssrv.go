package gowssrv

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ServerConfig struct {
	Origin string // * - all (no check)
}

type Server struct {
	address        string
	OnConnected    func(*Client)
	OnMessage      func(*Client, string)
	OnError        func(*Client, error)
	OnDisconnected func(*Client)
}

func NewServer(address, endpoint string, config ServerConfig) *Server {
	s := &Server{
		address: address,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		log.Print("origin: " + r.RemoteAddr)
		if config.Origin == "*" {
			return true
		}
		if config.Origin == r.RemoteAddr {
			return true
		}
		return false
	}

	http.HandleFunc(endpoint, s.handler)

	return s
}

func (s *Server) Serve() {
	http.ListenAndServe(s.address, nil)
}

func (s *Server) handler(rw http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	clientID, _ := uuid.GenerateUUID()
	client := &Client{
		connection: connection,
		ID:         clientID,
	}

	// Send client id to the jus connected client on startup
	connection.WriteMessage(websocket.TextMessage, []byte("client_id="+client.ID))
	s.OnConnected(client)
	for {
		_, msg, err := client.connection.ReadMessage()

		if err != nil {
			if cerr, ok := err.(*websocket.CloseError); ok {
				log.Printf("client %s is OnDisconnected | err: %s", client.ID, cerr.Error())
				s.OnDisconnected(client)
			} else {
				log.Printf("client %s is OnError | err: %s", client.ID, cerr.Error())
				s.OnError(client, err)
			}
			return
		}
		s.OnMessage(client, string(msg))
	}
}
