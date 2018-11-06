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

var mode=flag.String("m","","Execute or file transfer")

type Client struct {
	Addr string
	Port string
	Mode string
}

func NewClient(addr,port,mode string,)*Client{
	client:=&Client{
		Addr:addr,
		Port:port,
		Mode:mode,
	}
	return client
}

func interruptConn(conn net.Conn){
		conn.Close()
		os.Exit(0)
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
	log.Println("fd: ",fd)
	return fd
}

func SendToServer(conn net.Conn,filename,data string){
	log.Println("Sending to server data")
	filedata:=filename+"/"+data
	log.Println("Sending: ",filedata)
	_,err:=conn.Write([]byte(filedata))
	if err!=nil{
		log.Fatal("Error in sending to server filename: ",err)
	}
	log.Println("Successfully sending data to server")
}

//return filename from the file path
func getFilename(filepath string)(filename string){
	f:=strings.Split(filepath,"/")
	filename=f[len(f)-1]
	log.Printf("Psth: %v, filename: %v",filepath,filename)
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
	log.Println("Sended to server mode type")
	//read from server
	buf := make([]byte, 1024)
	k, er := conn.Read(buf)
	if er != nil || k == 0 {
		fmt.Println("Error reading from server: ", er)
		os.Exit(1)
	}
	fmt.Print(string(buf))
	var input string
	for {
		if c.Mode == "file" {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Input path to file to upload to server: ")
			scanner.Scan()
			filepath := scanner.Text()
			if filepath=="exit"{
				interruptConn(conn)
			}
			f,err:=os.Open(filepath)
			if err!=nil{
				log.Println("Error in opening file: ",err)
				fmt.Println(err)
				return
			}
			filename:=getFilename(filepath)
			fileData:=ReadFile(f)
			SendToServer(conn,filename,fileData)
			f.Close()
			log.Println("OK")
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Input: ")
			scanner.Scan()
			input = scanner.Text()
			if input == "exit" {
				interruptConn(conn)
			}
			//log.Println("sending to server: ",input)
			n, e := conn.Write([]byte(input))
			if e != nil || n == 0 {
				fmt.Println("Error:", err)
			}
			//log.Println("write to conn: ", input)
			//server
			fmt.Println("Server answer:")
			buff := make([]byte, 1024)
			k, er = conn.Read(buff)
			if er != nil || k == 0 {
				fmt.Println("Error reading from server: ", er)
				os.Exit(1)
			}
			//log.Println("get from server: ",string(buff))
			fmt.Println(string(buff))
			buff = buff[:0]
		}
	}
	conn.Close()
}

func main(){
	flag.Parse()
	c:=NewClient("","8080",*mode)
	log.Println("mode: ",c.Mode)
	c.startClient()
}