package main

import (
	"Driver-go/elevio"
	"flag"
	"fmt"
	"strconv"

	"HEISPROSJEKT/communication"
	//"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/hardware"
	"Network-go/network/bcast"
	"Network-go/network/peers"
)

func main() {
	fmt.Println("Jeg nekter å kommentere ut fmt hver gang jeg skal debugge")
	id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	numFloors := 4
	peerPort := 30004
	bcastPort := 30400

	//init functions
	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)
	hardwareChannels := hardware.InitElevatorHardware()
	channels := communication.InitNetworkChannels()
	//elevatorObject := elevatorHardware.InitializeElevatorObject(strconv.Itoa(*id))

	go peers.Receiver(peerPort, channels.PeerUpdateChl)
	go peers.Transmitter(peerPort, strconv.Itoa(*id), channels.PeerTxEnableCh)

	go bcast.Transmitter(bcastPort, channels.BcastOutgoingMessagesChannel)
	go bcast.Receiver(bcastPort, channels.BcastIncomingMessagesChannel)

	go communication.UpdatePeerList(channels)

	go elevio.PollButtons(hardwareChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(hardwareChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(hardwareChannels.PollObstructionChannel)
	go elevio.PollStopButton(hardwareChannels.PollStopButtonChannel)

	go hardware.RunElevatorHardware(strconv.Itoa(*id), hardwareChannels)

	go communication.BroadcastElevatorWorldView(strconv.Itoa(*id), channels.BcastOutgoingMessagesChannel, hardwareChannels.ElevatorObjectChannel)
	go communication.RecieveBroadcastfWorldViewfFromPeer(channels.BcastIncomingMessagesChannel)

	select {}
}
