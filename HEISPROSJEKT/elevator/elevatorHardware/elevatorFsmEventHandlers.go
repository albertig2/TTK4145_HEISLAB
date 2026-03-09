package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func InitElevatorBetweenFloors(elevator *elevatorConfig.Elevator) {
	ElevatorMotorDirection(elevatorConfig.Down)
	elevator.Direction = elevatorConfig.Down
	elevator.Behavior = elevatorConfig.Moving
}

func InitializeElevator(id string) elevatorConfig.Elevator {

	config := elevatorConfig.Config{
		DoorOpenDuration_s: 3 * time.Second,
	}

	elevator := elevatorConfig.Elevator{
		OwnId:     id,
		Floor:     elevio.GetFloor(),
		Direction: elevatorConfig.Stop,
		Requests:  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
		Behavior:  elevatorConfig.Idle,
		Config:    config,
	}
	return elevator
}

func InitElevatorHardwareChannels() elevatorConfig.ElevatorHardwareChannelsStruckt {

	hardwareChannels := elevatorConfig.ElevatorHardwareChannelsStruckt{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		FloorSensorChannel:      make(chan int),
		DoorOpenChannel:         make(chan bool),
		MotorDirectionChannel:   make(chan elevatorConfig.Direction),
		ElevatorObjectChannel:   make(chan elevatorConfig.Elevator),
	}
	return hardwareChannels
}

// small initialisation sequence to put the elevator in a known behavior
// small initialisation sequence to put elevator in a known state
func InitElevatorHardware(elevator *elevatorConfig.Elevator) {

	TurnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)
	InitElevatorBetweenFloors(elevator)

}

func updateMotorDirection(motorDirection chan elevatorConfig.Direction) {
	for {
		d := <-motorDirection

		println("Motordirection:", d)

		elevio.SetMotorDirection(elevio.MotorDirection(d))
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

func HandleOnFloorArrival(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, newFloor int, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Println("Elevator arrived at:", newFloor)
	debuggingHelpers.Elevator_print(*elevator)

	elevator.Floor = newFloor
	ElevatorFloorIndicatorLight(elevator.Floor)

	switch elevator.Behavior {
	case elevatorConfig.Moving:
		if RequestsShouldStop(*elevator) {

			ElevatorMotorDirection(elevatorConfig.Stop)

			//MotorDirectionChannel <- elevio.MD_Stop
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)
			//Elevator_doorLight(true)
			//*elevator = RequestsClearAtCurrentFloor(*elevator)
			//timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
			//SetAllLights(*elevator)
			//elevator.Behavior = elevatorConfig.DoorOpen
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state after HandleOnFloorArrival:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func OpenDoor(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, timeOpenSeconds time.Duration) {

	if !isBetweenFloors() {

		fmt.Println("Door Open")

		elevio.SetDoorOpenLamp(true)
		*elevator = RequestsClearAtCurrentFloor(*elevator)

		doorTimer.Stop()
		doorTimer.Reset(timeOpenSeconds)

		SetAllLights(*elevator)
		elevator.Behavior = elevatorConfig.DoorOpen

	} else {
		fmt.Println("Door was attempted opend in between floors")
		//Should probably trigger some sort of error handelig or somthing here
	}
}

func OnDoorTimeout(elevator *elevatorConfig.Elevator, doorTimer *time.Timer) {
	fmt.Println("Door timeout")
	debuggingHelpers.Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		pair := Requests_chooseDirection(*elevator)
		elevator.Direction = pair.Direction
		elevator.Behavior = pair.behavior

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			// timer.Timer_start(elevator.Config.DoorOpenDuration_s)
			// *elevator = RequestsClearAtCurrentFloor(*elevator)
			// SetAllLights(*elevator)
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			ElevatorDoorLight(false)
			ElevatorMotorDirection(elevator.Direction)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state after  OnDoorTimeout:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func HandleRequestButtonPress(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, btn_floor int, btn_type elevatorConfig.Button, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Printf("\n\n%s(%d, %s)\n", "Recieved the following order:", btn_floor, elevatorConfig.ButtonToString(btn_type))
	debuggingHelpers.Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		if RequestsShouldClearImmediately(*elevator, btn_floor, btn_type) {
			doorTimer.Stop()
			doorTimer.Reset(elevatorConfig.DOOR_OPEN_DURATION_S)
			//timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
		} else {
			elevator.Requests[btn_floor][btn_type] = true
		}

	case elevatorConfig.Moving:
		elevator.Requests[btn_floor][btn_type] = true

	case elevatorConfig.Idle:
		elevator.Requests[btn_floor][btn_type] = true
		pair := Requests_chooseDirection(*elevator)
		elevator.Direction = pair.Direction
		elevator.Behavior = pair.behavior

		switch pair.behavior {

		case elevatorConfig.DoorOpen:
			// ElevatorDoorLight(true)
			// timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
			// *elevator = RequestsClearAtCurrentFloor(*elevator)
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)

		case elevatorConfig.Moving:
			ElevatorMotorDirection(elevator.Direction)

		case elevatorConfig.Idle:
			// Do nothing
		}
	}

	SetAllLights(*elevator)

	fmt.Printf("\nNew state after HandleRequestButtonPres:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func HandleStopButtonActivated(stopActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {

	if stopActivated {
		//MotorDirectionChannel <- elevatorConfig.Stop

		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			//If stop is triggerd while at a floor, the door is opend and keep open until stop is reset + 3 seconds more
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)
			doorTimer.Stop()

		case elevatorConfig.Moving:
			ElevatorMotorDirection(elevatorConfig.Stop)

		default:

		}
		fmt.Println("Stop was activated")
		TurnOffAllOrderLights()
		elevio.SetStopLamp(true)

	} else {
		fmt.Println("Stop was reset")
		elevio.SetStopLamp(false)

		//MotorDirectionChannel <- elevatorConfig.Stop

		switch elevator.Behavior {

		case elevatorConfig.Moving, elevatorConfig.Idle:
			InitElevatorBetweenFloors(elevator)

		case elevatorConfig.DoorOpen:
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S) // keep door open for 3 more sek

		default:

		}

	}

}

// func HandleStopButtonreset(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {
// 	fmt.Println("Stop was activated")
// 	TurnOffAllOrderLights()
// 	elevio.SetStopLamp(true)

// 	//MotorDirectionChannel <- elevatorConfig.Stop

// 	switch elevator.Behavior {

// 	case elevatorConfig.Moving, elevatorConfig.Idle:
// 		InitElevatorBetweenFloors(elevator)

// 	default:

// 	}

// }

func HandleObstructionActivated(obstructionActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {

	if obstructionActivated {
		fmt.Println("Obstruction was activated")
		//ElevatorMotorDirection(elevatorConfig.Stop)
		//MotorDirectionChannel <- elevatorConfig.Stop

		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			doorTimer.Stop()

		default:
			//Do nothing if the door is not open
		}

	} else {
		fmt.Println("Obstruction was reset")
		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)

		default:
			//Do nothing if the door is not open
		}

	}

}
