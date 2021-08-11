package gowssrv

import "github.com/gorilla/websocket"

type Client struct {
	connection *websocket.Conn
	ID         string
}

func (c *Client) SendMessage(msg string) error {
	return c.connection.WriteMessage(websocket.TextMessage, []byte(msg))
}
