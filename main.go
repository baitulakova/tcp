package main

import (
	"net"
	"os"
	"fmt"
	"log"
)

func handleConnection(conn net.Conn){
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	conn.Write(buf)
}

func main(){
	listener,err:=net.Listen("tcp", ":8080")
	if err!=nil{
		log.Fatal("Error listening: ",err)
		os.Exit(1)
	}

	for{
		conn,err:=listener.Accept()
		if err!=nil{
			panic(err)
		}
		conn.Write([]byte("TCP Socket\n"))
		go handleConnection(conn)
	}
}
