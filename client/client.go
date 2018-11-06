package main

import (
	"net"
	"log"
	"fmt"
	"os"
	"bufio"
	"flag"
	"io"
	"strings"
)

var mode=flag.String("m","","Uploading file or routine mode")

type Client struct {
	Addr string
	Port string
	Mode string
}

func NewClient(addr,port,mode string)*Client{
	client:=&Client{
		Addr:addr,
		Port:port,
		Mode:mode,
	}
	return client
}

func interruptConn(conn net.Conn){
		conn.Close()
		fmt.Println("Connection closed")
}

func ReadFile(file *os.File)string{
	data:=make([]byte,1024)
	fd:=""
	for{
		n,err:=file.Read(data)
		if err==io.EOF{
			break
		}else if err!=nil{
			log.Fatalf("Error reading file - %v: %v",file.Name(),err)
		}
		fd=string(data[:n])
	}
	return fd
}

func SendToServer(conn net.Conn,filename,data string){
	log.Println("Sending to server data")
	filedata:=filename+"/"+data
	_,err:=conn.Write([]byte(filedata))
	if err!=nil{
		log.Fatal("Error in sending to server filename: ",err)
	}
	log.Println("Successfully sending data to server")
}

//getFilename return filename from the file path
func getFilename(filepath string)(filename string){
	f:=strings.Split(filepath,"/")
	filename=f[len(f)-1]
	log.Printf("Path: %v, filename: %v",filepath,filename)
	return filename
}

func (c *Client) startClient(){
	conn,err:=net.Dial("tcp",c.Addr+":"+c.Port)
	if err!=nil{
		log.Fatal("Error connecting to the server: ",err)
	}
	//send mode type
	_,e:=conn.Write([]byte(c.Mode))
	if e!=nil{
		log.Fatal("Error in sending mode type to server")
	}
	//read from server
	buf := make([]byte, 1024)
	k, er := conn.Read(buf)
	if er != nil || k == 0 {
		fmt.Println("Error reading from server: ", er)
		interruptConn(conn)
	}
	var input string
	for {
		if c.Mode == "file" {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Input path to file to upload to server: ")
			scanner.Scan()
			filepath := scanner.Text()
			if filepath=="exit"{
				break
			}
			f,err:=os.Open(filepath)
			if err!=nil{
				log.Println("Error in opening file: ",err)
				break
			}
			filename:=getFilename(filepath)
			fileData:=ReadFile(f)
			SendToServer(conn,filename,fileData)
			f.Close()
		} else if c.Mode=="r"{
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Input: ")
			scanner.Scan()
			input = scanner.Text()
			if input == "exit" {
				break
			}
			n, e := conn.Write([]byte(input))
			if e != nil || n == 0 {
				fmt.Println("Error in sending input text to server:", err)
				break
			}
			//server
			fmt.Println("Server's answer is:")
			buff := make([]byte, 1024)
			k, er = conn.Read(buff)
			if er != nil || k == 0 {
				fmt.Println("Error reading from server: ", er)
				break
			}
			fmt.Println(string(buff[:k]))
			buff = buff[:0]
		}else{
			fmt.Println("You indicated incorrect flag. Please use -m=file or -m=r")
			break
		}
	}
	interruptConn(conn)
}

func main(){
	flag.Parse()
	c:=NewClient("","8080",*mode)
	c.startClient()
}