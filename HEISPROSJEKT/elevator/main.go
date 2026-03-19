package main

import (
	"Driver-go/elevio"
	"flag"
	"fmt"
	"strconv"

	// "HEISPROSJEKT/communication"

	elevatorConfig "HEISPROSJEKT/elevator_config"
	elevatorController "HEISPROSJEKT/elevator_controller"
	orderProtocol "HEISPROSJEKT/order_protocol"
	"HEISPROSJEKT/synchronization"

	//"HEISPROSJEKT/elevatorHardware"
	//"HEISPROSJEKT/hardware"
	"Network-go/network/bcast"
	"Network-go/network/peers"
)

//note: Det skjer noen ganger at heisen kommer out of bounds, burde vi legge på faktiske hardware
//sikkert som hindrer heisen i å få til dette

//der er flere initer i kontroller, burde noen slås sammen?

//sto er broken i master, er den det her?

func main() {
	fmt.Println("Jeg nekter å kommentere ut fmt hver gang jeg skal debugge")
	id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	elevio.Init("localhost:"+strconv.Itoa(*port), elevatorConfig.NumberOfFloors)
	ControllerChannels := elevatorController.InitializeControllerChannels()
	orderChannels := orderProtocol.InitializeOrderChannels()
	synchronizationChannels := synchronization.InitializeSynchronizationChannels()

	go peers.Receiver(elevatorConfig.PeerUpdatePort, synchronizationChannels.PeerUpdateChannel)
	go peers.Transmitter(elevatorConfig.PeerUpdatePort, strconv.Itoa(*id), synchronizationChannels.PeerTransmitEnableChannel)

	go bcast.Transmitter(elevatorConfig.BroadcastPort, synchronizationChannels.BroadcastOutgoingMessagesChannel)
	go bcast.Receiver(elevatorConfig.BroadcastPort, synchronizationChannels.BroadcastIncomingMessagesChannel)

	go elevio.PollButtons(ControllerChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(ControllerChannels.PollFloorSensorChannel)
	go elevio.PollObstructionSwitch(ControllerChannels.PollObstructionChannel)
	go elevio.PollStopButton(ControllerChannels.PollStopButtonChannel)

	go elevatorController.LocalElevatorController(strconv.Itoa(*id), ControllerChannels, synchronizationChannels, orderChannels)

	go synchronization.SynchronizeElevators(synchronizationChannels, ControllerChannels, strconv.Itoa(*id))

	go orderProtocol.ManageAndDistributeOrders(strconv.Itoa(*id), orderChannels, synchronizationChannels, ControllerChannels)

	select {}
}
