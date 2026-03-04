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


		/*
		saves the motor direction before it stops.
		Right now it is used to continue moving in the same direction after stop or obstruction was activated
		this is most likely just a feature needed for testing the code for now. Should probably find a more logical way to deal with this if nescessary
		*/
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

func OpenDoor(doorTimer *time.Timer, timeOpenSeconds time.Duration){

	if(! isBetweenFloors()){

		currentElevatorState = DOOROPEN
		fmt.Println("Door Open")
		doorTimer.Reset(timeOpenSeconds)
		elevio.SetDoorOpenLamp(true)
	} else {
		fmt.Println("Door was attempted opend in between floors")
		//Should probably trigger some sort of error handelig or somthing here
	}
}

// func EnforceHardwareFloorBounderies(currentDirection elevio.MotorDirection, ) elevio.MotorDirection{

// }

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

			OpenDoor(doorTimer, 3*time.Second)
		case stopActivated := <-hardwareChannels.PollStopButtonChannel:
			if stopActivated {
				fmt.Println("Stop was activated")

				TurnOffAllOrderLights()
				elevio.SetStopLamp(true)
				hardwareChannels.MotorDirectionChannel <- elevio.MD_Stop

				if !isBetweenFloors() {
					OpenDoor(doorTimer, 3*time.Second)
				}
					
			} else {
				fmt.Println("Stop was Reset")

				elevio.SetStopLamp(false)
				if (isBetweenFloors()){
					hardwareChannels.MotorDirectionChannel <- _lastKnownDirection
				}

			}
		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
			if obstructionActivated {
				if currentElevatorState == DOOROPEN {
					fmt.Println("Obstruction was activated")
					hardwareChannels.MotorDirectionChannel <- elevio.MD_Stop
					doorTimer.Stop()
				} else {
					//Ignore input from obstruction if between floors/door is closed
				}

			} else {
				fmt.Println("Obstruction was Reset")
				OpenDoor(doorTimer, 3*time.Second) //keeps door open for 3 more seconds after obstruction was cleard
			}

		case <-doorTimer.C:
			elevio.SetDoorOpenLamp(false)
			fmt.Println("Door Closed")
			hardwareChannels.MotorDirectionChannel <- _nextMotorDirection
			currentElevatorState = MOVING
		}
	}
}
