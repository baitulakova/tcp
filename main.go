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

func (c *Client) Addr() string{
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

//GetData read data from connection
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
	if command=="file"{
		log.Println("Client want to upload file")
			Uploadfile(c.conn)
			log.Println("Successfully uploaded file")
	}else if command=="r"{
		log.Println("Client chose routine mode")
		for {
			input:=GetData(c.conn)
			_, e := c.conn.Write([]byte(input))
			if e != nil {
				fmt.Println("Error sending to client: ", e)
			}
			log.Println("Sending to client: ", string(input))
			input = input[:0]
		}
	} else{
		log.Println("Indicated incorrect flag. Flag was: ",command)
	}
	c.Close()
}

//CreateStorage create storage for files which will be uploaded to server
func CreateStorage() (path string){
	userHome:=os.Getenv("HOME")
	path = userHome+"/TCPServerStorage/"
	err:=os.MkdirAll(path,os.ModePerm)
	if err!=nil{
		fmt.Println("Error creating file storage: ",err)
	}
	return path
}

func Uploadfile(conn net.Conn){
	log.Println("File uploading started")
	defer conn.Close()
	filedata:=GetData(conn)
	fd:=strings.Split(filedata,"/")
	if len(fd)==0{
		log.Println("Error filename and data are empty.")
		conn.Close()
	}
	filename:=fd[0]
	data:=fd[1]
	log.Println("Successfully got data from client")
	log.Println("Creating file storage")
	fs:=CreateStorage()
	log.Println("Succesfully created file storage at: ",fs)
	log.Println("Creating file")
	file, err := os.Create(fs+filename)
	if err != nil {
		log.Fatal("Error in creation file: ",err)
	}
	log.Println("Successfully created file ",filename)
	defer file.Close()
	_,err=file.WriteString(data)
	if err!=nil{
		log.Fatal("Error in writing data to file: ",err)
	}
	log.Println("Uploaded file: ",filename)
}

func main(){
	s:=NewServer("","8080")
	s.Listen()
}
