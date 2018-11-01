package main

import (
	"net"
	"log"
	"fmt"
	"os"
)

type Server struct {
	address string
	port string
}

func NewServer(addr,p string) *Server{
	log.Printf("Creating new server with address %v and port %v",addr,p)
	server:=&Server{
		address:addr,
		port:p,
	}
	return server
}

func (s *Server) Listen(){
	log.Println("Listening to server: ",s.address+":"+s.port)
	listener,err:=net.Listen("tcp", s.address+":"+s.port)
	if err!=nil{
		log.Fatal("Error listening: ",err)
		os.Exit(1)
	}
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			log.Println("Error accepting connections: ",err)
			panic(err)
		}
		client:=&Client{
			conn:conn,
		}
		go client.handleConnection()
	}
}

type Client struct{
	conn net.Conn
}

func (c *Client) Close(){
	c.conn.Close()
}

func (c *Client) handleConnection(){
	log.Println("Listening for connection on: ",c.conn.RemoteAddr().String())
	defer c.Close()
	buf := make([]byte, 1024)
	_, err := c.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	c.conn.Write(buf)
}

func main(){
	s:=NewServer("","8080")
	s.Listen()
}
