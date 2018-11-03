package main

import (
	"net"
	"log"
	"fmt"
	"os"
	"bufio"
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
	//io.Copy(os.Stdout, conn)
	buf := make([]byte, 1024)
	k, er := conn.Read(buf)
	if er != nil || k == 0 {
		fmt.Println("Error reading from server: ", er)
		os.Exit(1)
	}
	fmt.Print(string(buf))
	var input string
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Input: ")
		scanner.Scan()
		input = scanner.Text()
		if input=="exit"{
			break
		}
		//log.Println("sending to server: ",input)
		n, e := conn.Write([]byte(input))
		if e != nil || n == 0 {
			fmt.Println("Error:", err)
		}
		//log.Println("write to conn: ", input)
		//server
		fmt.Println("Server answer:")
		buf = make([]byte, 1024)
		k, er = conn.Read(buf)
		if er != nil || k == 0 {
			fmt.Println("Error reading from server: ", er)
			os.Exit(1)
		}
		//log.Println("get from server: ",string(buf))
		fmt.Println(string(buf))
	}
	conn.Close()
}

func main(){
	c:=NewClient("","8080")
	c.startClient()
}