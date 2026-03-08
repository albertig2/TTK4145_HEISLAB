package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/timer"
	"fmt"
)

// func SetAllLights(elevator elevatorConfig.Elevator) {
// 	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
// 		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
// 			Elevator_requestButtonLight(floor, elevatorConfig.Button(btn), elevator.Requests[floor][btn])
// 		}
// 	}
// }

func Fsm_onInitBetweenFloors(elevator *elevatorConfig.Elevator) {
	ElevatorMotorDirection(elevatorConfig.Down)
	elevator.Direction = elevatorConfig.Down
	elevator.Behavior = elevatorConfig.Moving
}

func Fsm_onRequestButtonPress(elevator *elevatorConfig.Elevator, btn_floor int, btn_type elevatorConfig.Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "Fsm_handleRequestButtonPress", btn_floor, elevatorConfig.ButtonToString(btn_type))
	Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		if Requests_shouldClearImmediately(*elevator, btn_floor, btn_type) {
			timer.Timer_start(elevator.Config.DoorOpenDuration_s)
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
			ElevatorDoorLight(true)
			timer.Timer_start(elevator.Config.DoorOpenDuration_s)
			*elevator = Requests_clearAtCurrentFloor(*elevator)

		case elevatorConfig.Moving:
			ElevatorMotorDirection(elevator.Direction)

		case elevatorConfig.Idle:
			// nothing
		}
	}

	SetAllLights(*elevator)

	fmt.Printf("\nNew state:\n")
	Elevator_print(*elevator)
}

func Fsm_onFloorArrival(elevator *elevatorConfig.Elevator, newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	Elevator_print(*elevator)

	elevator.Floor = newFloor
	ElevatorFloorIndicator(elevator.Floor)

	switch elevator.Behavior {
	case elevatorConfig.Moving:
		if Requests_shouldStop(*elevator) {
			ElevatorMotorDirection(elevatorConfig.Stop)
			ElevatorDoorLight(true)
			*elevator = Requests_clearAtCurrentFloor(*elevator)
			timer.Timer_start(elevator.Config.DoorOpenDuration_s)
			SetAllLights(*elevator)
			elevator.Behavior = elevatorConfig.DoorOpen
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*elevator)
}

func fsm_onDoorTimeout(elevator *elevatorConfig.Elevator) {
	fmt.Printf("\n\n%s()\n", "fsm_onDoorTimeout")
	Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		pair := Requests_chooseDirection(*elevator)
		elevator.Direction = pair.Direction
		elevator.Behavior = pair.behavior

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			timer.Timer_start(elevator.Config.DoorOpenDuration_s)
			*elevator = Requests_clearAtCurrentFloor(*elevator)
			SetAllLights(*elevator)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			ElevatorDoorLight(false)
			ElevatorMotorDirection(elevator.Direction)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*elevator)
}
