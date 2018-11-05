package main

import (
	"net"
	"log"
	"os"
	"fmt"
	"bufio"
	"io"
	"flag"
)

var(
	command = flag.String("c","","choose command")
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
	flag.Parse()
	log.Println("c: ",*command)
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
		go client.handleConnection(*command)
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

func (c *Client) handleConnection(command string){
	log.Println("Listening for connection on: ", c.Addr())
	c.SendString("Hello " + c.Addr() + "\n")
	if command=="file"{
		err:=fileTransfer(c.conn)
		if err!=nil{
			log.Fatal("Error in fileTransfer: ", err)
			os.Exit(1)
		}
	}else {
		for {
			input := make([]byte, 1024)
			n, err := c.conn.Read(input)
			log.Println("get from client: ", string(input))
			if err != nil || n == 0 {
				fmt.Println("Error reading:", err.Error())
				os.Exit(1)
			}
			if string(input) == "exit" {
				c.Close()
				os.Exit(1)
			}
			_, e := c.conn.Write(input)
			if e != nil {
				fmt.Println("error sending to client: ", e)
			}
			log.Println("sending to client: ", string(input))
			input = input[:0]
		}
	}
	c.conn.Close()
}

func fileTransfer(conn net.Conn) error {
	log.Println("In file transfer")
	defer conn.Close()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Input path to file: ")
	scanner.Scan()
	input := scanner.Text()
	f,errOpenFile:=os.Open(input)
	if errOpenFile!=nil{
		log.Fatal("File not found. No such file.")
		os.Exit(1)
	}
	file, err := os.Create(f.Name())
	if err != nil {
		log.Println("Error in create")
		return err
	}
	_, err = io.Copy(file, f)
	log.Println("End of transfer")
	return err
}

func main(){
	s:=NewServer("","8080")
	s.Listen()
}
