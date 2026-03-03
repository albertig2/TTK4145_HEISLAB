package main

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/Hardware"
	"flag"
	"strconv"



	
	// "Net"
	"fmt"
)

func main() {
	// id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	numFloors := 4

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)

	// go peers.Receiver(65004, peerUpdateChl)
	// go peers.Transmitter(65004, strconv.Itoa(*id), peerRecieveEnableChl)

	

	var d elevio.MotorDirection = elevio.MD_Up
	
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_mdir := make (chan elevio.MotorDirection)
	

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go Hardware.HardwareSafetyFeatures(drv_obstr, drv_stop, drv_mdir)
	go Hardware.MotorDriection(drv_mdir)

	drv_mdir <- d

	for {
		select {
			case a := <-drv_buttons:
				fmt.Printf("%+v\n", a)
				elevio.SetButtonLamp(a.Button, a.Floor, true)

			case a := <-drv_floors:
				fmt.Printf("%+v\n", a)
				if a == numFloors-1 {
					d = elevio.MD_Down
				} else if a == 0 {
					d = elevio.MD_Up
				}
				drv_mdir <- d
		}
	}
	
}
