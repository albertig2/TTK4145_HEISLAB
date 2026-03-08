package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
)

type DirectionBehaviorPair struct {
	Direction elevatorConfig.Direction
	behavior  elevatorConfig.Behavior
}

// static i C -> privat i Go (liten forbokstav)
func Requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < elevatorConfig.N_FLOORS; f++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e Elevator) bool {
	for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

// “API” lik headeren: behold navnet
func Requests_chooseDirection(e Elevator) DirectionBehaviorPair {
	switch e.Direction {
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
	switch e.Direction {
	case elevatorConfig.Down:
		return e.Requests[e.Floor][elevatorConfig.HallDown] ||
			e.Requests[e.Floor][elevatorConfig.Cab] ||
			!Requests_below(e)

	case elevatorConfig.Up:
		return e.Requests[e.Floor][elevatorConfig.HallUp] ||
			e.Requests[e.Floor][elevatorConfig.Cab] ||
			!Requests_above(e)

	case elevatorConfig.Stop:
		fallthrough
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_Floor int, btn_type elevatorConfig.Button) bool {
	return e.Floor == btn_Floor &&
		((e.Direction == elevatorConfig.Up && btn_type == elevatorConfig.HallUp) ||
			(e.Direction == elevatorConfig.Down && btn_type == elevatorConfig.HallDown) ||
			e.Direction == elevatorConfig.Stop ||
			btn_type == elevatorConfig.Cab)
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.Requests[e.Floor][elevatorConfig.Cab] = false

	switch e.Direction {
	case elevatorConfig.Up:
		if !Requests_above(e) && !e.Requests[e.Floor][elevatorConfig.HallUp] {
			e.Requests[e.Floor][elevatorConfig.HallDown] = false
		}
		e.Requests[e.Floor][elevatorConfig.HallUp] = false

	case elevatorConfig.Down:
		if !Requests_below(e) && !e.Requests[e.Floor][elevatorConfig.HallDown] {
			e.Requests[e.Floor][elevatorConfig.HallUp] = false
		}
		e.Requests[e.Floor][elevatorConfig.HallDown] = false

	case elevatorConfig.Stop:
		fallthrough
	default:
		e.Requests[e.Floor][elevatorConfig.HallUp] = false
		e.Requests[e.Floor][elevatorConfig.HallDown] = false
	}

	return e
}
