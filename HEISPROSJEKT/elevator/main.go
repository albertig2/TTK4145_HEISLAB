package main

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/Hardware"
	"flag"
	"strconv"
	// "Net"
	// "fmt"
)

func main() {
	// id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	numFloors := 4

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)

	// go peers.Receiver(65004, peerUpdateChl)
	// go peers.Transmitter(65004, strconv.Itoa(*id), peerRecieveEnableChl)

	hardwareChannels := Hardware.InitElevatorHardware()

	go elevio.PollButtons(hardwareChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(hardwareChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(hardwareChannels.PollObstructionChannel)
	go elevio.PollStopButton(hardwareChannels.PollStopButtonChannel)


	go Hardware.RunElevatorHardware(hardwareChannels)


	select {}

}
