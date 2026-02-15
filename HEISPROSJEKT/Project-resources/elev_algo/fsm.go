package elevator

import "fmt"

func setAllLights(es Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevator_requestButtonLight(floor, btn, es.requests[floor][btn])
		}
	}
}

func fsm_onInitBetweenFloors(e *Elevator) {
	elevator_motorDirection(D_Down)
	e.dirn = D_Down
	e.behaviour = EB_Moving
}

func fsm_onRequestButtonPress(e *Elevator, btn_floor int, btn_type Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btn_floor, elevator_buttonToString(btn_type))
	elevator_print(*e)

	switch e.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(*e, btn_floor, btn_type) {
			timer_start(e.config.doorOpenDuration_s)
		} else {
			e.requests[btn_floor][btn_type] = true
		}

	case EB_Moving:
		e.requests[btn_floor][btn_type] = true

	case EB_Idle:
		e.requests[btn_floor][btn_type] = true
		pair := requests_chooseDirection(*e)
		e.dirn = pair.dirn
		e.behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			elevator_doorLight(true)
			timer_start(e.config.doorOpenDuration_s)
			*e = requests_clearAtCurrentFloor(*e)

		case EB_Moving:
			elevator_motorDirection(e.dirn)

		case EB_Idle:
			// nothing
		}
	}

	setAllLights(*e)

	fmt.Printf("\nNew state:\n")
	elevator_print(*e)
}

func fsm_onFloorArrival(e *Elevator, newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	elevator_print(*e)

	e.floor = newFloor
	elevator_floorIndicator(e.floor)

	switch e.behaviour {
	case EB_Moving:
		if requests_shouldStop(*e) {
			elevator_motorDirection(D_Stop)
			elevator_doorLight(true)
			*e = requests_clearAtCurrentFloor(*e)
			timer_start(e.config.doorOpenDuration_s)
			setAllLights(*e)
			e.behaviour = EB_DoorOpen
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	elevator_print(*e)
}

func fsm_onDoorTimeout(e *Elevator) {
	fmt.Printf("\n\n%s()\n", "fsm_onDoorTimeout")
	elevator_print(*e)

	switch e.behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(*e)
		e.dirn = pair.dirn
		e.behaviour = pair.behaviour

		switch e.behaviour {
		case EB_DoorOpen:
			timer_start(e.config.doorOpenDuration_s)
			*e = requests_clearAtCurrentFloor(*e)
			setAllLights(*e)

		case EB_Moving, EB_Idle:
			elevator_doorLight(false)
			elevator_motorDirection(e.dirn)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	elevator_print(*e)
}
