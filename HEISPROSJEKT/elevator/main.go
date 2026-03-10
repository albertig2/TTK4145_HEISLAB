package main

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/communication"
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/hardware"
	"Network-go/network/bcast"
	"Network-go/network/peers"
	"flag"
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("Jeg nekter å kommentere ut fmt hver gang jeg skal debugge")
	id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	system := elevatorHardware.ElevatorSystem{}
	elevatorHardware.Initialize(&system, strconv.Itoa(*id))
	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorHardware.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorHardware.OrderStatus)
	receivedWorldView := make(chan string)

	go orderProtocol.orderRutine(&system, HallRequestsForAllElevators, CabRequestsForAllElevators, receivedWorldView)

	numFloors := 4
	peerPort := 30004
	bcastPort := 30400

	//init functions
	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)
	hardwareChannels := hardware.InitElevatorHardware()
	channels := communication.InitNetworkChannels()

	go peers.Receiver(peerPort, channels.PeerUpdateChl)
	go peers.Transmitter(peerPort, strconv.Itoa(*id), channels.PeerTxEnableCh)

	go bcast.Transmitter(bcastPort, channels.BcastOutgoingMessagesChannel)
	go bcast.Receiver(bcastPort, channels.BcastIncomingMessagesChannel)

	go communication.BroadcastElevatorWorldView(strconv.Itoa(*id), channels.BcastOutgoingMessagesChannel, hardwareChannels.ElevatorStateChannel)
	go communication.RecieveBroadcastfWorldViewfFromPeer(channels.BcastIncomingMessagesChannel)

	go communication.UpdatePeerList(channels)

	go elevio.PollButtons(hardwareChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(hardwareChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(hardwareChannels.PollObstructionChannel)
	go elevio.PollStopButton(hardwareChannels.PollStopButtonChannel)

	go hardware.RunElevatorHardware(hardwareChannels)

	select {}
}
