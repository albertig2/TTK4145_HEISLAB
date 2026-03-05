package main

import (
	"Driver-go/elevio"
	"flag"
	"strconv"

	//"Network-go/network/peers"
	"HEISPROSJEKT/communication"

	"fmt"
	"HEISPROSJEKT/Hardware"

)


func main() {
	id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	numFloors := 4
	peerPort := 65004

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)

	hardwareChannels := Hardware.InitElevatorHardware()

	channels := communication.InitNetworkChannels()

	communication.StartPeerNetworking(peerPort, *id, channels)

	
	
	for {
		
		communication.UpdatePeerList(channels)
		fmt.Printf("Peers: %q\n", communication.GetAlivePeersList())
		fmt.Printf("Dead: %q\n", communication.GetDeadPeersList())
	}
	

	go elevio.PollButtons(hardwareChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(hardwareChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(hardwareChannels.PollObstructionChannel)
	go elevio.PollStopButton(hardwareChannels.PollStopButtonChannel)


	go Hardware.RunElevatorHardware(hardwareChannels)


	select {}
} 

	


