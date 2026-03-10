package hardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

var _stopActivated bool = false

const _doorOpenTime = 3 * time.Second

var _doorIsOpen bool = false

var currentElevatorBehavior Behavior = elevatorConfig.Idle

var _nextMotorDirection elevatorConfig.Direction = elevatorConfig.Up
var _lastKnownDirection elevatorConfig.Direction = elevatorConfig.Stop

func updateMotorDirection(motorDirection chan elevatorConfig.Direction) {
	for {
		d := <-motorDirection
		println("Motordirection:", d)

		/*
			saves the motor direction before it stops.
			Right now it is used to continue moving in the same direction after stop or obstruction was activated
			this is most likely just a feature needed for testing the code for now. Should probably find a more logical way to deal with this if nescessary
		*/
		if d != elevatorConfig.Stop {
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
	MotorDirectionChannel   chan elevatorConfig.Direction
	ElevatorBehaviorChannel chan ElevatorBehavior
}

func InitElevatorHardware() ElevatorHardwareChannelsStruckt {

	hardwareChannels := ElevatorHardwareChannelsStruckt{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		FloorSensorChannel:      make(chan int),
		DoorOpenChannel:         make(chan bool),
		MotorDirectionChannel:   make(chan elevatorConfig.Direction),
		ElevatorBehaviorChannel: make(chan ElevatorBehavior),
	}

	//small initialisation sequence to put the elevator in a known behavior
	TurnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)

	var initialDirection elevatorConfig.Direction

	if elevio.GetFloor() != elevatorConfig.N_FLOORS-1 {
		initialDirection = elevatorConfig.Up
	} else {
		initialDirection = elevatorConfig.Down
	}
	elevio.SetMotorDirection(initialDirection)

	fmt.Printf("Motordirection was set to %+v when running init \n", initialDirection)

	return hardwareChannels
}

func OpenDoor(doorTimer *time.Timer, timeOpenSeconds time.Duration) {

	if !isBetweenFloors() {

		currentElevatorBehavior = elevatorConfig.DoorOpen
		fmt.Println("Door Open")
		doorTimer.Reset(timeOpenSeconds)
		elevio.SetDoorOpenLamp(true)
	} else {
		fmt.Println("Door was attempted opend in between floors")
		//Should probably trigger some sort of error handelig or somthing here
	}
}

// func EnforceHardwareFloorBounderies(currentDirection elevatorConfig.Direction, ) elevatorConfig.Direction{

// }

func TriggerStopButtonSideEffects(doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Println("Stop was activated")
	TurnOffAllOrderLights()
	elevio.SetStopLamp(true)
	MotorDirectionChannel <- elevatorConfig.Stop

	if !isBetweenFloors() { //If stop is triggerd while at a floor, the door is opend and keep open until stop is reset + 3 seconds more
		OpenDoor(doorTimer, 3*time.Second)
		doorTimer.Stop()
	}

}

func TriggerObstructionSideEffects(doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {

	fmt.Println("Obstruction was activated")
	MotorDirectionChannel <- elevatorConfig.Stop
	doorTimer.Stop()

}

func RunElevatorHardware(hardwareChannels ElevatorHardwareChannelsStruckt) {

	doorTimer := time.NewTimer(_doorOpenTime)
	doorTimer.Stop()

	go updateMotorDirection(hardwareChannels.MotorDirectionChannel)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:
			hardwareChannels.MotorDirectionChannel <- elevatorConfig.Stop
			fmt.Printf("Elevator arrived at floor %+v\n", floor)
			if floor == elevatorConfig.N_FLOORS-1 {
				_nextMotorDirection = elevatorConfig.Down
			} else if floor == 0 {
				_nextMotorDirection = elevatorConfig.Up
			}

			OpenDoor(doorTimer, 3*time.Second)
		case stopActivated := <-hardwareChannels.PollStopButtonChannel:
			if stopActivated {

				TriggerStopButtonSideEffects(doorTimer, hardwareChannels.MotorDirectionChannel)

				/*
					TurnOffAllOrderLights()
					elevio.SetStopLamp(true)
					hardwareChannels.MotorDirectionChannel <- elevatorConfig.Stop

					if !isBetweenFloors() {
						OpenDoor(doorTimer, 3*time.Second)
					}

				*/

			} else {
				fmt.Println("Stop was Reset")

				elevio.SetStopLamp(false)
				if isBetweenFloors() {
					hardwareChannels.MotorDirectionChannel <- _lastKnownDirection
				} else {
					OpenDoor(doorTimer, 3*time.Second) //keeps door open for 3 more seconds after obstruction was cleard

				}

			}
		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
			if obstructionActivated {
				if currentElevatorBehavior == elevatorConfig.DoorOpen {
					// fmt.Println("Obstruction was activated")
					// hardwareChannels.MotorDirectionChannel <- elevatorConfig.Stop
					// doorTimer.Stop()
					TriggerObstructionSideEffects(doorTimer, hardwareChannels.MotorDirectionChannel)

				} else {
					//Ignore input from obstruction if between door is closed
				}

			} else {
				fmt.Println("Obstruction was Reset")
				OpenDoor(doorTimer, 3*time.Second) //keeps door open for 3 more seconds after obstruction was cleard
			}

		case <-doorTimer.C:
			elevio.SetDoorOpenLamp(false)
			fmt.Println("Door Closed")
			hardwareChannels.MotorDirectionChannel <- _nextMotorDirection
			currentElevatorBehavior = elevatorConfig.Moving
			hardwareChannels.ElevatorBehaviorChannel <- currentElevatorBehavior
		}
	}
}
