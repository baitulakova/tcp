package main

import (
	"net"
	"log"
	"os"
	"fmt"
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
	defer listener.Close()
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

func (c *Client) Addr()string{
	return c.conn.RemoteAddr().String()
}

func (c *Client) Close(){
	c.conn.Close()
	log.Printf("Connection with %v closed",c.Addr())
}

func (c *Client) SendString(message string){
	c.conn.Write([]byte(message))
}

func (c *Client) handleConnection(){
	log.Println("Listening for connection on: ", c.Addr())
	c.SendString("Hello " + c.Addr() + "\n")
	log.Println("send hello")
	for {
		input := make([]byte, 1024)
		n, err := c.conn.Read(input)
		if string(input)=="exit"{
			break
		}
		log.Println("get from client: ",string(input))
		if err != nil || n == 0 {
			fmt.Println("Error reading:", err.Error())
			os.Exit(1)
		}
		c.SendString("You entered: ")
		_, e := c.conn.Write(input)
		if e != nil {
			fmt.Println("error sending to client: ", e)
		}
		log.Println("sending to client: ",string(input))
		input=input[:0]
	}
	c.conn.Close()
}

func main(){
	s:=NewServer("","8080")
	s.Listen()
}
