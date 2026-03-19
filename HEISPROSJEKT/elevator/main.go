package main

import (
	"Driver-go/elevio"
	elevatorConfig "HEISPROSJEKT/elevator_config"
	elevatorController "HEISPROSJEKT/elevator_controller"
	orderProtocol "HEISPROSJEKT/order_protocol"
	"HEISPROSJEKT/synchronization"
	"Network-go/network/bcast"
	"Network-go/network/peers"
	"flag"
	"strconv"
)

func main() {
	ownId := flag.String("id", "1", "Input id")
	elevatorServerPort := flag.Int("port", 15657, "Input port")
	flag.Parse()

	elevio.Init("localhost:"+strconv.Itoa(*elevatorServerPort), elevatorConfig.NumberOfFloors)
	ControllerChannels := elevatorController.InitializeControllerChannels()
	orderChannels := orderProtocol.InitializeOrderChannels()
	synchronizationChannels := synchronization.InitializeSynchronizationChannels()

	go peers.Receiver(elevatorConfig.PeerUpdatePort, synchronizationChannels.PeerUpdateChannel)
	go peers.Transmitter(elevatorConfig.PeerUpdatePort, *ownId, synchronizationChannels.PeerTransmitEnableChannel)

	go bcast.Transmitter(elevatorConfig.BroadcastPort, synchronizationChannels.BroadcastOutgoingMessagesChannel)
	go bcast.Receiver(elevatorConfig.BroadcastPort, synchronizationChannels.BroadcastIncomingMessagesChannel)

	go elevio.PollButtons(ControllerChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(ControllerChannels.PollFloorSensorChannel)
	go elevio.PollObstructionSwitch(ControllerChannels.PollObstructionChannel)
	go elevio.PollStopButton(ControllerChannels.PollStopButtonChannel)

	go elevatorController.LocalElevatorController(*ownId, ControllerChannels, synchronizationChannels, orderChannels)

	go synchronization.SynchronizeElevators(synchronizationChannels, ControllerChannels, *ownId)

	go orderProtocol.ManageAndDistributeOrders(*ownId, orderChannels, synchronizationChannels, ControllerChannels)

	select {}
}
