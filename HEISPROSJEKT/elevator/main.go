package main

import (
	"Driver-go/elevio"
	"flag"
	"strconv"

	"Network-go/network/peers"
	"HEISPROSJEKT/communication"

	"fmt"
)

//func testTransmit(peerUpdateChl chan<- peers.PeerUpdate) {

//	peers.Receiver(65004, peerUpdateChl)

//}

func main() {
	id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	numFloors := 4

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)

	peerUpdateChl := make(chan peers.PeerUpdate)
	peerRecieveEnableChl := make(chan bool)

	go peers.Receiver(65004, peerUpdateChl)
	go peers.Transmitter(65004, strconv.Itoa(*id), peerRecieveEnableChl)
	
	
	for {
		
		a := <-peerUpdateChl
		fmt.Printf("Peers: %q\n", a.Peers)
		
	}
} 
