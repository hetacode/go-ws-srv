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
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: just for beginning phase
}

type Server struct {
	address        string
	OnConnected    func(Client)
	OnError        func(Client, error)
	OnDisconnected func(Client)
}

func NewServer(address, endpoint string) *Server {
	http.HandleFunc(endpoint, func(rw http.ResponseWriter, r *http.Request) {
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
		for {
			_, msg, err := connection.ReadMessage()
			if err != nil {
				return
			}
			panic(msg)
		}
	})
	s := &Server{
		address: address,
	}

	return s
}

func (s *Server) Serve() {
	http.ListenAndServe(s.address, nil)
}
