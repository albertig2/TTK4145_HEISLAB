package main

import (
	"Driver-go/elevio"
	"flag"
	"fmt"
	"strconv"

	// "HEISPROSJEKT/communication"

	elevatorController "HEISPROSJEKT/elevator_controller"
	"HEISPROSJEKT/orderProtocol"
	"HEISPROSJEKT/synchronization"

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

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)
	ControllerChannels := elevatorController.InitializeControllerChannels()
	orderChannels := orderProtocol.InitializeOrderChannels()
	synchronizationChannels := synchronization.InitializeSynchronizationChannels()

	go peers.Receiver(peerPort, synchronizationChannels.PeerUpdateChannel)
	go peers.Transmitter(peerPort, strconv.Itoa(*id), synchronizationChannels.PeerTxEnableChannel)

	go bcast.Transmitter(bcastPort, synchronizationChannels.BcastOutgoingMessagesChannel)
	go bcast.Receiver(bcastPort, synchronizationChannels.BcastIncomingMessagesChannel)

	go elevio.PollButtons(ControllerChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(ControllerChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(ControllerChannels.PollObstructionChannel)
	go elevio.PollStopButton(ControllerChannels.PollStopButtonChannel)

	go elevatorController.LocalElevatorController(strconv.Itoa(*id), ControllerChannels, synchronizationChannels, orderChannels)

	go synchronization.SynchronizeElevators(ControllerChannels.LocalElevatorChannel, synchronizationChannels,ControllerChannels, strconv.Itoa(*id))

	go orderProtocol.ManageAndDistributeOrders(strconv.Itoa(*id), orderChannels, synchronizationChannels, ControllerChannels)

	select {}
}
