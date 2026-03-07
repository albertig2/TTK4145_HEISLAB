package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/timer"
	"fmt"
)

func SetAllLights(es Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			Elevator_requestButtonLight(floor, Button(btn), es.requests[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors(e *Elevator) {
	Elevator_motorDirection(elevatorConfig.Down)
	e.direction = elevatorConfig.Down
	e.behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(e *Elevator, btn_floor int, btn_type Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btn_floor, Elevator_buttonToString(btn_type))
	Elevator_print(*e)

	switch e.behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(*e, btn_floor, btn_type) {
			timer.Timer_start(e.config.doorOpenDuration_s)
		} else {
			e.requests[btn_floor][btn_type] = true
		}

	case EB_Moving:
		e.requests[btn_floor][btn_type] = true

	case EB_Idle:
		e.requests[btn_floor][btn_type] = true
		pair := requests_chooseDirection(*e)
		e.direction = pair.direction
		e.behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			Elevator_doorLight(true)
			timer.Timer_start(e.config.doorOpenDuration_s)
			*e = Requests_clearAtCurrentFloor(*e)

		case EB_Moving:
			Elevator_motorDirection(e.direction)

		case EB_Idle:
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

	e.floor = newFloor
	Elevator_floorIndicator(e.floor)

	switch e.behaviour {
	case EB_Moving:
		if Requests_shouldStop(*e) {
			Elevator_motorDirection(elevatorConfig.Stop)
			Elevator_doorLight(true)
			*e = Requests_clearAtCurrentFloor(*e)
			timer.Timer_start(e.config.doorOpenDuration_s)
			SetAllLights(*e)
			e.behaviour = EB_DoorOpen
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

	switch e.behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(*e)
		e.direction = pair.direction
		e.behaviour = pair.behaviour

		switch e.behaviour {
		case EB_DoorOpen:
			timer.Timer_start(e.config.doorOpenDuration_s)
			*e = Requests_clearAtCurrentFloor(*e)
			SetAllLights(*e)

		case EB_Moving, EB_Idle:
			Elevator_doorLight(false)
			Elevator_motorDirection(e.direction)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*e)
}
