package Hardware

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

var _numFloors int = 4


var _stopActivated bool = false

const _doorOpenTime = 3 * time.Second

var _doorIsOpen bool = false

type elevatorState int

const (
	IDLE     elevatorState = 0
	MOVING   elevatorState = 1
	DOOROPEN elevatorState = 2
)

var currentElevatorState elevatorState = IDLE
var _nextMotorDirection elevio.MotorDirection = elevio.MD_Up
var _lastKnownDirection elevio.MotorDirection = elevio.MD_Stop

func MotorDriection(motorDirection chan elevio.MotorDirection) {
	for {
		d := <-motorDirection
		println("Motordirection:", d)

		if d != elevio.MD_Stop {

			_lastKnownDirection = d
		}

		elevio.SetMotorDirection(d)
	}
}


// func OpenDoor(openDoorChannel chan bool, motorDirection chan elevio.MotorDirection) {

// 	_doorTimer := time.NewTimer(_doorOpenTime)
// 	_doorTimer.Stop()

// 	for {
// 		select {
// 		case openDoor := <-openDoorChannel:
// 			fmt.Println("Door open was triggerd")
// 			if openDoor && !(isBetweenFloors()) {
// 				fmt.Println("Door open loop was triggerd")
// 				motorDirection <- elevio.MD_Stop
// 				_doorTimer.Reset(_doorOpenTime)
// 				_doorIsOpen = true
// 				elevio.SetDoorOpenLamp(true)
// 			} else {
// 				//do nothing if false
// 			}

// 		case <-_doorTimer.C:
// 			fmt.Println("Timeout was triggerd")
// 			_doorIsOpen = false
// 			elevio.SetDoorOpenLamp(false)
// 			motorDirection <- (_lastKnownDirection)
// 		}
// 	}
// }

// func manegOrderLight(){

// }

func isBetweenFloors() bool {
	currentFloor := elevio.GetFloor()
	if currentFloor != -1 {
		return false
	} else {
		return true
	}

}

type ElevatorHardwareChannelsStruckt struct {
	PollOrderButtonsChannel chan elevio.ButtonEvent
	PollObstructionChannel chan bool
	PollStopButtonChannel  chan bool
	FloorSensorChannel           chan int
	DoorOpenChannel        chan bool
	MotorDirectionChannel  chan elevio.MotorDirection
	ElevatorStateChannel   chan elevatorState
}

func InitElevatorHaredwareChannels() ElevatorHardwareChannelsStruckt {

	hardwareChannels := ElevatorHardwareChannelsStruckt{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel: make(chan bool),
		PollStopButtonChannel:  make(chan bool),
		FloorSensorChannel:           make(chan int),
		DoorOpenChannel:        make(chan bool),
		MotorDirectionChannel:  make(chan elevio.MotorDirection),
		ElevatorStateChannel:   make(chan elevatorState),
	}

	return hardwareChannels
}

func RunElevatorHardware(hardwareChannels ElevatorHardwareChannelsStruckt) {

	doorTimer := time.NewTimer(_doorOpenTime)
	doorTimer.Stop()

	go MotorDriection(hardwareChannels.MotorDirectionChannel)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:
			hardwareChannels.MotorDirectionChannel <- elevio.MD_Stop
			fmt.Printf("Elevator arrived at floor %+v\n", floor)
			if floor == _numFloors-1 {
				_nextMotorDirection = elevio.MD_Down
			} else if floor == 0 {
				_nextMotorDirection = elevio.MD_Up
			} 
			fmt.Println("Next direction is now ", _nextMotorDirection)

			

			currentElevatorState = DOOROPEN
			fmt.Println("Door open was triggerd")
			doorTimer.Reset(_doorOpenTime)
			elevio.SetDoorOpenLamp(true)

		// case stopActivated := <-hardwareChannels.pollStopButtonChannel:

		// case obstructionActivated := <-hardwareChannels.pollObstructionChannel:

		// case doorOpen := <-hardwareChannels.doorOpenChannel:

		// case motorDirection := <-hardwareChannels.motorDirectionChannel:

		case <- doorTimer.C:
			fmt.Println("Timeout was triggerd")
			_doorIsOpen = false
			elevio.SetDoorOpenLamp(false)
			fmt.Println("Door closed was triggerd")
			hardwareChannels.MotorDirectionChannel <- _nextMotorDirection
			// hardwareChannels.ElevatorStateChannel <- MOVING
		}
	}
}

		// case state := <-hardwareChannels.ElevatorStateChannel:
		// 	switch state {

		// 	case IDLE:
		// 		currentElevatorState = IDLE
		// 		hardwareChannels.MotorDirectionChannel <- elevio.MD_Stop

		// 	case MOVING:
		// 		currentElevatorState = MOVING
		// 		hardwareChannels.MotorDirectionChannel <- _nextMotorDirection
		// 	case DOOROPEN:
		// 		currentElevatorState = DOOROPEN

		// 		hardwareChannels.MotorDirectionChannel <- elevio.MD_Stop
		// 		fmt.Println("Door open was triggerd")
		// 		doorTimer.Reset(_doorOpenTime)
		// 		_doorIsOpen = true
		// 		elevio.SetDoorOpenLamp(true)
		// 	}

// func HardwareSafetyFeatures(pollObstructionChannel chan bool, pollStopButtonChannel chan bool, openDoorChannel chan bool, motorDirection chan elevio.MotorDirection) {
// 	for {
// 		select {

// 		case obstructionActivated := <-pollObstructionChannel:

// 			// fmt. Println(obstructionActivated)
// 			currentFloor := elevio.GetFloor()

// 			if obstructionActivated {
// 				if currentFloor != -1 {
// 					fmt.Println("Obstruction was activated")
// 					motorDirection <- elevio.MD_Stop
// 					// openDoorChannel<- true

// 				} else {
// 					//Ignore input from obstruction if between floors/door is closed
// 				}

// 			} else {
// 				fmt.Println("Obstruction was Reset")
// 				motorDirection <- _lastKnownDirection
// 			}

// 		case stopActivated := <-pollStopButtonChannel:
// 			fmt.Println(stopActivated)

// 			if stopActivated {
// 				fmt.Println("Stop was activated")

// 				elevio.SetStopLamp(true)
// 				motorDirection <- elevio.MD_Stop
// 				TurnOffAllLights()

// 			} else {
// 				fmt.Println("Stop was Reset")

// 				elevio.SetStopLamp(false)
// 				motorDirection <- _lastKnownDirection

// 			}

// 		}
// 	}
// }
