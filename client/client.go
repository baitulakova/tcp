package main

import (
	"net"
	"log"
	"fmt"
	"os"
)

type Client struct {
	Addr string
	Port string
}

func NewClient(a string,p string)*Client{
	client:=&Client{
		Addr:a,
		Port:p,
	}
	return client
}

func (c *Client) startClient(){
	conn,err:=net.Dial("tcp",c.Addr+":"+c.Port)
	if err!=nil{
		log.Fatal("Error connecting to the server: ",err)
	}
	defer conn.Close()
	//io.Copy(os.Stdout, conn)
	buf := make([]byte, 1024)
	k, er := conn.Read(buf)
	if er != nil || k == 0 {
		fmt.Println("Error reading from server: ", er)
		os.Exit(1)
	}
	fmt.Println(string(buf))
	var input string
	fmt.Print("Input: ")
	fmt.Scanln(&input)
	n, e := conn.Write([]byte(input))
	if e != nil || n == 0 {
		fmt.Println("Error:", err.Error())
	}

	//server
	fmt.Println("Server answer:")
	buf = make([]byte, 1024)
	k, er = conn.Read(buf)
	if er != nil || k == 0 {
		fmt.Println("Error reading from server: ", er)
		os.Exit(1)
	}
	fmt.Println(string(buf))
}

func main(){
	c:=NewClient("","8080")
	c.startClient()
}