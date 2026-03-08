package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/timer"
	"fmt"
)

func SetAllLights(es Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			Elevator_requestButtonLight(floor, elevatorConfig.Button(btn), es.Requests[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors(e *Elevator) {
	Elevator_motorDirection(elevatorConfig.Down)
	e.Direction = elevatorConfig.Down
	e.Behavior = elevatorConfig.Moving
}

func Fsm_onRequestButtonPress(e *Elevator, btn_floor int, btn_type elevatorConfig.Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btn_floor, elevatorConfig.ButtonToString(btn_type))
	Elevator_print(*e)

	switch e.Behavior {
	case elevatorConfig.DoorOpen:
		if Requests_shouldClearImmediately(*e, btn_floor, btn_type) {
			timer.Timer_start(e.Config.DoorOpenDuration_s)
		} else {
			e.Requests[btn_floor][btn_type] = true
		}

	case elevatorConfig.Moving:
		e.Requests[btn_floor][btn_type] = true

	case elevatorConfig.Idle:
		e.Requests[btn_floor][btn_type] = true
		pair := Requests_chooseDirection(*e)
		e.Direction = pair.Direction
		e.Behavior = pair.behavior

		switch pair.behavior {
		case elevatorConfig.DoorOpen:
			Elevator_doorLight(true)
			timer.Timer_start(e.Config.DoorOpenDuration_s)
			*e = Requests_clearAtCurrentFloor(*e)

		case elevatorConfig.Moving:
			Elevator_motorDirection(e.Direction)

		case elevatorConfig.Idle:
			// nothing
		}
	}

	SetAllLights(*e)

	fmt.Printf("\nNew state:\n")
	Elevator_print(*e)
}

func Fsm_onFloorArrival(e *Elevator, newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	Elevator_print(*e)

	e.Floor = newFloor
	Elevator_floorIndicator(e.Floor)

	switch e.Behavior {
	case elevatorConfig.Moving:
		if Requests_shouldStop(*e) {
			Elevator_motorDirection(elevatorConfig.Stop)
			Elevator_doorLight(true)
			*e = Requests_clearAtCurrentFloor(*e)
			timer.Timer_start(e.Config.DoorOpenDuration_s)
			SetAllLights(*e)
			e.Behavior = elevatorConfig.DoorOpen
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*e)
}

func fsm_onDoorTimeout(e *Elevator) {
	fmt.Printf("\n\n%s()\n", "fsm_onDoorTimeout")
	Elevator_print(*e)

	switch e.Behavior {
	case elevatorConfig.DoorOpen:
		pair := Requests_chooseDirection(*e)
		e.Direction = pair.Direction
		e.Behavior = pair.behavior

		switch e.Behavior {
		case elevatorConfig.DoorOpen:
			timer.Timer_start(e.Config.DoorOpenDuration_s)
			*e = Requests_clearAtCurrentFloor(*e)
			SetAllLights(*e)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			Elevator_doorLight(false)
			Elevator_motorDirection(e.Direction)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*e)
}
