package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/timer"
	"fmt"
)

func SetAllLights(es Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			Elevator_requestButtonLight(floor, Button(btn), es.Requests[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors(e *Elevator) {
	Elevator_motorDirection(elevio.MD_Down)
	e.Dirn = elevio.MD_Down
	e.Behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(e *Elevator, btn_floor int, btn_type Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btn_floor, Elevator_buttonToString(btn_type))
	Elevator_print(*e)

	switch e.Behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(*e, btn_floor, btn_type) {
			timer.Timer_start(e.Config.DoorOpenDuration_s)
		} else {
			e.Requests[btn_floor][btn_type] = true
		}

	case EB_Moving:
		e.Requests[btn_floor][btn_type] = true

	case EB_Idle:
		e.Requests[btn_floor][btn_type] = true
		pair := requests_chooseDirection(*e)
		e.Dirn = pair.dirn
		e.Behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			Elevator_doorLight(true)
			timer.Timer_start(e.Config.DoorOpenDuration_s)
			*e = Requests_clearAtCurrentFloor(*e)

		case EB_Moving:
			Elevator_motorDirection(e.Dirn)

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

	e.Floor = newFloor
	Elevator_floorIndicator(e.Floor)

	switch e.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(*e) {
			Elevator_motorDirection(D_Stop)
			Elevator_doorLight(true)
			*e = Requests_clearAtCurrentFloor(*e)
			timer.Timer_start(e.Config.DoorOpenDuration_s)
			SetAllLights(*e)
			e.Behaviour = EB_DoorOpen
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

	switch e.Behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(*e)
		e.Dirn = pair.dirn
		e.Behaviour = pair.behaviour

		switch e.Behaviour {
		case EB_DoorOpen:
			timer.Timer_start(e.Config.DoorOpenDuration_s)
			*e = Requests_clearAtCurrentFloor(*e)
			SetAllLights(*e)

		case EB_Moving, EB_Idle:
			Elevator_doorLight(false)
			Elevator_motorDirection(e.Dirn)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*e)
}
