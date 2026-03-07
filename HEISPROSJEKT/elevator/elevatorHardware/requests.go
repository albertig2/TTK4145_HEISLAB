package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
)

type DirectionBehaviourPair struct {
	direction elevatorConfig.Direction
	behaviour ElevatorBehaviour
}

// static i C -> privat i Go (liten forbokstav)
func Requests_above(e Elevator) bool {
	for f := e.floor + 1; f < elevatorConfig.N_FLOORS; f++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e Elevator) bool {
	for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return true
		}
	}
	return false
}

// “API” lik headeren: behold navnet
func requests_chooseDirection(e Elevator) DirectionBehaviourPair {
	switch e.direction {
	case elevatorConfig.Up:
		if Requests_above(e) {
			return DirectionBehaviourPair{elevatorConfig.Up, EB_Moving}
		}
		if Requests_here(e) {
			return DirectionBehaviourPair{elevatorConfig.Down, EB_DoorOpen}
		}
		if Requests_below(e) {
			return DirectionBehaviourPair{elevatorConfig.Down, EB_Moving}
		}
		return DirectionBehaviourPair{elevatorConfig.Stop, EB_Idle}

	case elevatorConfig.Down:
		if Requests_below(e) {
			return DirectionBehaviourPair{elevatorConfig.Down, EB_Moving}
		}
		if Requests_here(e) {
			return DirectionBehaviourPair{elevatorConfig.Up, EB_DoorOpen}
		}
		if Requests_above(e) {
			return DirectionBehaviourPair{elevatorConfig.Up, EB_Moving}
		}
		return DirectionBehaviourPair{elevatorConfig.Stop, EB_Idle}

	case elevatorConfig.Stop: // samme kommentar som i C
		if Requests_here(e) {
			return DirectionBehaviourPair{elevatorConfig.Stop, EB_DoorOpen}
		}
		if Requests_above(e) {
			return DirectionBehaviourPair{elevatorConfig.Up, EB_Moving}
		}
		if Requests_below(e) {
			return DirectionBehaviourPair{elevatorConfig.Down, EB_Moving}
		}
		return DirectionBehaviourPair{elevatorConfig.Stop, EB_Idle}

	default:
		return DirectionBehaviourPair{elevatorConfig.Stop, EB_Idle}
	}
}

func Requests_shouldStop(e Elevator) bool {
	switch e.direction {
	case elevatorConfig.Down:
		return e.requests[e.floor][elevatorConfig.HallDown] ||
			e.requests[e.floor][elevatorConfig.Cab] ||
			!Requests_below(e)

	case elevatorConfig.Up:
		return e.requests[e.floor][elevatorConfig.HallUp] ||
			e.requests[e.floor][elevatorConfig.Cab] ||
			!Requests_above(e)

	case elevatorConfig.Stop:
		fallthrough
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevatorConfig.Button) bool {
	return e.floor == btn_floor &&
		((e.direction == elevatorConfig.Up && btn_type == elevatorConfig.HallUp) ||
			(e.direction == elevatorConfig.Down && btn_type == elevatorConfig.HallDown) ||
			e.direction == elevatorConfig.Stop ||
			btn_type == elevatorConfig.Cab)
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.requests[e.floor][elevatorConfig.Cab] = false

	switch e.direction {
	case elevatorConfig.Up:
		if !Requests_above(e) && !e.requests[e.floor][elevatorConfig.HallUp] {
			e.requests[e.floor][elevatorConfig.HallDown] = false
		}
		e.requests[e.floor][elevatorConfig.HallUp] = false

	case elevatorConfig.Down:
		if !Requests_below(e) && !e.requests[e.floor][elevatorConfig.HallDown] {
			e.requests[e.floor][elevatorConfig.HallUp] = false
		}
		e.requests[e.floor][elevatorConfig.HallDown] = false

	case elevatorConfig.Stop:
		fallthrough
	default:
		e.requests[e.floor][elevatorConfig.HallUp] = false
		e.requests[e.floor][elevatorConfig.HallDown] = false
	}

	return e
}
