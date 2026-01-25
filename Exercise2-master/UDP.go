package main

import (
	"fmt"
	"net"
	"time"
)

func recieve() {
	//reciever socket

	addr := "0.0.0.0:20004" //Listening address, should maybe bee 2000 + 4?

	server, err := net.ResolveUDPAddr("udp", addr) //resolve connection to server wi

	if err != nil {
		fmt.Println(err)
	}

	connection, _ := net.ListenUDP("udp", server) //listen to server

	buffer := make([]byte, 1024)

	defer connection.Close()

	for {

		n, senderAddr, _ := connection.ReadFromUDP(buffer)

		fmt.Println("Message:", string(buffer[0:n]), " recieved from ", senderAddr) //prints message from server and the server address

		time.Sleep(50 * time.Millisecond)

	}
}

func transmitt() {
	//laddr, _ := net.ResolveUDPAddr("udp", ":port") //resolve connection to server wi

	raddr, _ := net.ResolveUDPAddr("udp", "10.100.23.11:20004") //resolve connection to server win

	conn, _ := net.DialUDP("udp", nil, raddr)

	defer conn.Close()

	for {

		message := "Hello world from station 4"

		data := []byte(message)

		conn.Write(data)

		time.Sleep(50 * time.Millisecond)
	}
}

func main() {

	//runtime.GOMAXPROCS(2)

	go recieve()
	go transmitt()

	select {}

	//sending socket
}
