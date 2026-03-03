package main

import (
	"Driver-go/elevio"
	"flag"
	"strconv"

	//"Network-go/network/peers"
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
	peerPort := 65004

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)


	channels := communication.InitNetworkChannels(*id, peerPort)


	
	
	for {
		
		communication.UpdatePeerList(channels)
		fmt.Printf("Peers: %q\n", communication.GetAlivePeersList())
		fmt.Printf("Dead: %q\n", communication.GetDeadPeersList())
	}
} 
