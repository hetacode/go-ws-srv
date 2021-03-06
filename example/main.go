package main

import (
	"log"

	gowssrv "github.com/hetacode/go-ws-srv"
)

func main() {
	s := gowssrv.NewServer(":8080", "/", gowssrv.ServerConfig{Origin: "*"})
	s.OnConnected = func(c *gowssrv.Client) {
		log.Printf("client %s OnConnected", c.ID)
	}
	s.OnMessage = func(c *gowssrv.Client, s string) {
		log.Printf("client %s msg: %s", c.ID, s)
	}
	s.OnDisconnected = func(c *gowssrv.Client) {
		log.Printf("client %s OnDisconnected", c.ID)
	}
	s.OnError = func(c *gowssrv.Client, e error) {
		log.Printf("client %s OnError %s", c.ID, e)
	}
	s.Serve()
}
