package main

import (
	"Driver-go/elevio"
	"flag"
	"fmt"
	"strconv"

	// "HEISPROSJEKT/communication"

	elevatorController "HEISPROSJEKT/elevator_controller"
	orderProtocol "HEISPROSJEKT/order_protocol"
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
	elevatorControllerChannels := elevatorController.InitializeControllerChannels()
	orderChannels := orderProtocol.InitializeOrderChannels()
	synchronizationChannels := synchronization.InitializeSynchronizationChannels()

	go peers.Receiver(peerPort, synchronizationChannels.PeerUpdateChannel)
	go peers.Transmitter(peerPort, strconv.Itoa(*id), synchronizationChannels.PeerTxEnableChannel)

	go bcast.Transmitter(bcastPort, synchronizationChannels.BcastOutgoingMessagesChannel)
	go bcast.Receiver(bcastPort, synchronizationChannels.BcastIncomingMessagesChannel)

	go elevio.PollButtons(elevatorControllerChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(elevatorControllerChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(elevatorControllerChannels.PollObstructionChannel)
	go elevio.PollStopButton(elevatorControllerChannels.PollStopButtonChannel)

	go elevatorController.LocalElevatorController(strconv.Itoa(*id), elevatorControllerChannels, synchronizationChannels, orderChannels)

	go synchronization.SynchronizeElevators(elevatorControllerChannels.LocalElevatorChannel, synchronizationChannels, strconv.Itoa(*id))

	go orderProtocol.ManageAndDistributeOrders(strconv.Itoa(*id), orderChannels, synchronizationChannels, elevatorControllerChannels)

	select {}
}
