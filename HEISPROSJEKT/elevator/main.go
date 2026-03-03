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

	//elevio.SetMotorDirection(d)

	// drv_buttons := make(chan elevio.ButtonEvent)
	// drv_floors := make(chan int)
	// drv_obstr := make(chan bool)
	// drv_stop := make(chan bool)
	// drv_mdir := make (chan elevio.MotorDirection)
	// drv_doorOpen := make(chan bool)

	hardwareChannels := Hardware.InitElevatorHardware()

	go elevio.PollButtons(hardwareChannels.PollOrderButtonsChannel)
	go elevio.PollFloorSensor(hardwareChannels.FloorSensorChannel)
	go elevio.PollObstructionSwitch(hardwareChannels.PollObstructionChannel)
	go elevio.PollStopButton(hardwareChannels.PollStopButtonChannel)
	// go Hardware.HardwareSafetyFeatures(drv_obstr, drv_stop, drv_doorOpen, drv_mdir)
	// go Hardware.MotorDriection(drv_mdir)
	// go Hardware.OpenDoor(drv_doorOpen, make(chan elevio.MotorDirection))

	go Hardware.RunElevatorHardware(hardwareChannels)

	// drv_mdir <- d
	select {}

	// for {
	// 	select {
	// 		case a := <-hardwareChannels.PollOrderButtonsChannel:
	// 			fmt.Printf("%+v\n", a)
	// 			elevio.SetButtonLamp(a.Button, a.Floor, true)

	// 		case a := <- hardwareChannels.FloorSensorChannel:
	// 			fmt.Printf("%+v\n", a)

	// 			if a == numFloors-1 {
	// 				d = elevio.MD_Down
	// 			} else if a == 0 {
	// 				d = elevio.MD_Up
	// 			}
	// 			elevio.SetMotorDirection(d)

	// 	}
	// }

}
