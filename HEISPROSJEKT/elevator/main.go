package main

import (
	"Driver-go/elevio"
	"flag"
	"strconv"

	"Network-go/network/peers"
	"HEISPROSJEKT/communication"

	"fmt"
)

//func testTransmit(peerUpdateChl chan<- peers.PeerUpdate) {

//	peers.Receiver(65004, peerUpdateChl)

//}

func main() {
	id := flag.Int("id", 1, "Input id")
	port := flag.Int("port", 15657, "Input port")
	flag.Parse()

	fmt.Println(communication.Hello())

	numFloors := 4

	elevio.Init("localhost:"+strconv.Itoa(*port), numFloors)

	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	peerUpdateChl := make(chan peers.PeerUpdate)
	peerRecieveEnableChl := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go peers.Receiver(65004, peerUpdateChl)
	go peers.Transmitter(65004, strconv.Itoa(*id), peerRecieveEnableChl)
	//peerRecieveEnableChl <- true
	
	

	for {
		//select {
		// case a := <-drv_buttons:
		// 	fmt.Printf("%+v\n", a)
		// 	elevio.SetButtonLamp(a.Button, a.Floor, true)

		// case a := <-drv_floors:
		// 	fmt.Printf("%+v\n", a)
		// 	if a == numFloors-1 {
		// 		d = elevio.MD_Down
		// 	} else if a == 0 {
		// 		d = elevio.MD_Up
		// 	}
		// 	elevio.SetMotorDirection(d)

		// case a := <-drv_obstr:
		// 	fmt.Printf("%+v\n", a)
		// 	if a {
		// 		elevio.SetMotorDirection(elevio.MD_Stop)
		// 	} else {
		// 		elevio.SetMotorDirection(d)
		// 	}

		// case a := <-drv_stop:
		// 	fmt.Printf("%+v\n", a)
		// 	for f := 0; f < numFloors; f++ {
		// 		for b := elevio.ButtonType(0); b < 3; b++ {
		// 			elevio.SetButtonLamp(b, f, false)
		// 		}
		// 	}
		//case a := <-drv_stop:
		//	fmt.Printf("%+v\n", a)
		//	for f := 0; f < numFloors; f++ {
		//		for b := elevio.ButtonType(0); b < 3; b++ {
		//			elevio.SetButtonLamp(b, f, false)
		//		}
		//	}
		//case

		a := <-peerUpdateChl
		fmt.Printf("Peers: %q\n", a.Peers)
		//}
	}
} //hei
