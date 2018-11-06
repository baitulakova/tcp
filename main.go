package main

import (
	"net"
	"log"
	"os"
	"fmt"
	"strings"
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

func (c *Client) ChooseMode()string{
	buf := make([]byte, 1024)
	n, er := c.conn.Read(buf)
	if er != nil{
		fmt.Println("Error reading from client mode type: ", er)
		os.Exit(1)
	}
	return string(buf[:n])
}

func (c *Client) SendString(message string){
	c.conn.Write([]byte(message))
}

//get data from client
func GetData(conn net.Conn)string{
	data:=make([]byte,1024)
	n,err:=conn.Read(data)
	if err!=nil{
		log.Fatal("Error reading data from client: ",err)
	}
	return string(data[:n])
}

func (c *Client) handleConnection(){
	log.Println("Listening for connection on: ", c.Addr())
	c.SendString("Hello " + c.Addr() + "\n")
	command:=c.ChooseMode()
	log.Println("Got mode type: ",command)
	if command=="file"{
		log.Println("Client want to transfer file")
		err:=fileTransfer(c.conn)
		if err!=nil{
			log.Fatal("Error in fileTransfer: ", err)
			os.Exit(1)
		}
	}else if command=="exe"{
		log.Println("Client want to execute")
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

func createStorage() (path string){
	userHome:=os.Getenv("HOME")
	path = userHome+"/TCPServerStorage/"
	err:=os.MkdirAll(path,os.ModePerm)
	if err!=nil{
		fmt.Println("Error creating file storage: ",err)
	}
	return path
}

func fileTransfer(conn net.Conn) error {
	log.Println("File uploading started")
	defer conn.Close()
	filedata:=GetData(conn)
	fd:=strings.Split(filedata,"/")
	filename:=fd[0]
	data:=fd[1]
	log.Println("Creating file storage")
	fs:=createStorage()
	log.Println("Succesfully created file storage at: ",fs)
	log.Println("Creating file")
	file, err := os.Create(fs+filename)
	if err != nil {
		log.Println("Error in create file strorage")
		return err
	}
	log.Println("Successfully created file ",filename)
	defer file.Close()
	log.Println("Successfully got data from client")
	file.WriteString(data)
	log.Println("Uploaded file: ",filename)
	return err
}

func main(){
	s:=NewServer("","8080")
	s.Listen()
}
