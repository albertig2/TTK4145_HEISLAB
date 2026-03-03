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

func updateMotorDirection(motorDirection chan elevio.MotorDirection) {
	for {
		d := <-motorDirection
		println("Motordirection:", d)

		if d != elevio.MD_Stop {
			_lastKnownDirection = d
		}

		elevio.SetMotorDirection(d)
	}
}

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
	PollObstructionChannel  chan bool
	PollStopButtonChannel   chan bool
	FloorSensorChannel      chan int
	DoorOpenChannel         chan bool
	MotorDirectionChannel   chan elevio.MotorDirection
	ElevatorStateChannel    chan elevatorState
}

func InitElevatorHardware() ElevatorHardwareChannelsStruckt {

	hardwareChannels := ElevatorHardwareChannelsStruckt{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		FloorSensorChannel:      make(chan int),
		DoorOpenChannel:         make(chan bool),
		MotorDirectionChannel:   make(chan elevio.MotorDirection),
		ElevatorStateChannel:    make(chan elevatorState),
	}

	//small initialisation sequence to put the elevator in a known state
	TurnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)

	var initialDirection elevio.MotorDirection = elevio.MD_Down

	if elevio.GetFloor() != _numFloors-1 {
		initialDirection = elevio.MD_Up
	} else {
		initialDirection = elevio.MD_Down
	}
	elevio.SetMotorDirection(initialDirection)

	fmt.Printf("Motordirection was set to %+v when running init \n", initialDirection)

	return hardwareChannels
}

func RunElevatorHardware(hardwareChannels ElevatorHardwareChannelsStruckt) {

	doorTimer := time.NewTimer(_doorOpenTime)
	doorTimer.Stop()

	go updateMotorDirection(hardwareChannels.MotorDirectionChannel)

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

		case <-doorTimer.C:
			fmt.Println("Timeout was triggerd")
			_doorIsOpen = false
			elevio.SetDoorOpenLamp(false)
			fmt.Println("Door closed was triggerd")
			hardwareChannels.MotorDirectionChannel <- _nextMotorDirection
			// hardwareChannels.ElevatorStateChannel <- MOVING
		}
	}
}
