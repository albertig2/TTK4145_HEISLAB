package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
)

type DirectionBehaviorPair struct {
	direction elevatorConfig.Direction
	behavior  elevatorConfig.Behavior
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
func requests_chooseDirection(e Elevator) DirectionBehaviorPair {
	switch e.direction {
	case elevatorConfig.Up:
		if Requests_above(e) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		if Requests_here(e) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.DoorOpen}
		}
		if Requests_below(e) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	case elevatorConfig.Down:
		if Requests_below(e) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		if Requests_here(e) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.DoorOpen}
		}
		if Requests_above(e) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	case elevatorConfig.Stop: // samme kommentar som i C
		if Requests_here(e) {
			return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.DoorOpen}
		}
		if Requests_above(e) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		if Requests_below(e) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	default:
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}
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
