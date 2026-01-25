package main

import (
	"fmt"
	"net"
	"time"
)

const mainServerIP = "10.100.23.11"
const mainServerPort = "33546"

const localIP = "10.100.23.14"
const localListenPort = "30069"

func tcp_send() {
	serverAddr, _ := net.ResolveTCPAddr("tcp", mainServerIP+":"+mainServerPort)
	connection, _ := net.DialTCP("tcp", nil, serverAddr)
	defer connection.Close()

	message := fmt.Sprintf("Connect to: %s:%s\x00", localIP, localListenPort)
	data := []byte(message)
	connection.Write(data)
}

func conn_handler(connection net.Conn) {
	for {
		buffer := make([]byte, 1024)
		connection.Write([]byte("Hello world\x00"))
		n, _ := connection.Read(buffer)
		fmt.Println("Recieved: ", string(buffer[:n]))
		time.Sleep(50 * time.Millisecond)
	}
}

func tcp_listen() {
	localAddr, _ := net.ResolveTCPAddr("tcp", localIP+":"+localListenPort)
	listener, _ := net.ListenTCP("tcp", localAddr)
	defer listener.Close()

	connection, _ := listener.AcceptTCP()
	defer connection.Close()

	conn_handler(connection)
}

func main() {
	go tcp_listen()

	time.Sleep(50 * time.Millisecond)

	go tcp_send()

	select {}
}
