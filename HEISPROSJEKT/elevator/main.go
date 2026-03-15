package main

import (
	"Driver-go/elevio"
	"flag"
	"fmt"
	"strconv"

	// "HEISPROSJEKT/communication"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/orderProtocol"
	"HEISPROSJEKT/synchronisation"

	//"HEISPROSJEKT/elevatorHardware"
	//"HEISPROSJEKT/hardware"
	"Network-go/network/bcast"
	"Network-go/network/peers"
)

//note to slef: DET ER TO FUNSKJONER SOM LYTTER PÅ PEERUPDATE CHL (eller kanskje dte ble skrevet om nå i kveld(fredag)
//dette er grunnen til at peerUpdat oppfører seg rart, og at den ikke printes fra sync elevators

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
	hardwareChannels := elevatorHardware.InitElevatorHardwareChannels()
	orderChannels := debuggingHelpers.InitializeOrderChannels()
	channels := synchronisation.InitNetworkChannels()
	//elevatorObject := elevatorHardware.InitializeElevatorObject(strconv.Itoa(*id))

	go peers.Receiver(peerPort, channels.PeerUpdateChannel)
	go peers.Transmitter(peerPort, strconv.Itoa(*id), channels.PeerTxEnableChannel)

	go bcast.Transmitter(bcastPort, channels.BcastOutgoingMessagesChannel)
	go bcast.Receiver(bcastPort, channels.BcastIncomingMessagesChannel)

	go synchronisation.UpdatePeerList(channels)

	go elevio.PollButtons(hardwareChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(hardwareChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(hardwareChannels.PollObstructionChannel)
	go elevio.PollStopButton(hardwareChannels.PollStopButtonChannel)

	//go debuggingHelpers.MimicOrderAssignerAndSynch(orderChannels)

	go elevatorHardware.RunElevatorFsm(strconv.Itoa(*id), hardwareChannels, channels, orderChannels)

	go synchronisation.SynchroniseElevators(hardwareChannels.ElevatorObjectChannel, channels, strconv.Itoa(*id))

	go orderProtocol.RunOrder(strconv.Itoa(*id), orderChannels, channels, hardwareChannels)

	// go communication.BroadcastElevatorWorldView(strconv.Itoa(*id), channels.BcastOutgoingMessagesChannel, hardwareChannels.ElevatorObjectChannel)
	// go communication.RecieveBroadcastfWorldViewfFromPeer(channels.BcastIncomingMessagesChannel)

	select {}
}
