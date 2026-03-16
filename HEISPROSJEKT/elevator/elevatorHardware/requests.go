package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
)

type DirectionBehaviorPair struct {
	Direction elevatorConfig.Direction
	behavior  elevatorConfig.Behavior
}

// static i C -> privat i Go (liten forbokstav)
func RequestsAbove(elevator elevatorConfig.Elevator) bool {
	for f := elevator.Floor + 1; f < elevatorConfig.N_FLOORS; f++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if elevator.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(elevator elevatorConfig.Elevator) bool {
	for f := 0; f < elevator.Floor; f++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if elevator.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestsHere(elevator elevatorConfig.Elevator) bool {
	for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
		if elevator.Requests[elevator.Floor][btn] {
			return true
		}
	}
	return false
}

// “API” lik headeren: behold navnet
func Requests_chooseDirection(elevator elevatorConfig.Elevator) DirectionBehaviorPair {
	switch elevator.Direction {
	case elevatorConfig.Up:
		if RequestsAbove(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		if RequestsHere(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.DoorOpen}
		}
		if RequestsBelow(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	case elevatorConfig.Down:
		if RequestsBelow(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		if RequestsHere(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.DoorOpen}
		}
		if RequestsAbove(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	case elevatorConfig.Stop: // samme kommentar som i C
		if RequestsHere(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.DoorOpen}
		}
		if RequestsAbove(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		if RequestsBelow(elevator) {
			return DirectionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	default:
		return DirectionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}
	}
}

func RequestsShouldStop(elevator elevatorConfig.Elevator) bool {
	switch elevator.Direction {
	case elevatorConfig.Down:
		return elevator.Requests[elevator.Floor][elevatorConfig.HallDown] ||
			elevator.Requests[elevator.Floor][elevatorConfig.Cab] ||
			!RequestsBelow(elevator)

	case elevatorConfig.Up:
		return elevator.Requests[elevator.Floor][elevatorConfig.HallUp] ||
			elevator.Requests[elevator.Floor][elevatorConfig.Cab] ||
			!RequestsAbove(elevator)

	case elevatorConfig.Stop:
		fallthrough
	default:
		return true
	}
}

func RequestsShouldClearImmediately(elevator elevatorConfig.Elevator, btn_Floor int, btn_type elevatorConfig.Button) bool {
	return elevator.Floor == btn_Floor &&
		((elevator.Direction == elevatorConfig.Up && btn_type == elevatorConfig.HallUp) ||
			(elevator.Direction == elevatorConfig.Down && btn_type == elevatorConfig.HallDown) ||
			elevator.Direction == elevatorConfig.Stop ||
			btn_type == elevatorConfig.Cab)
}

//fix so that the order that is cleard here, also is sendt on the new order serviced channel
//note to order handeling: this function clears order regardless of there actually was an oorder there. This means that the order 
//state machine has to be able to handle instances where it recieves a CleardOrder message for orders in other states than assigned 
func RequestsClearAtCurrentFloor(e elevatorConfig.Elevator, ServicedOrderChannel chan elevatorConfig.ButtonEvent) elevatorConfig.Elevator {

	e.Requests[e.Floor][elevatorConfig.Cab] = false
	ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.Cab} //clear light after the order after network is pinged
	ElevatorRequestButtonLight(e.Floor, elevatorConfig.Cab, false)

	switch e.Direction {
	case elevatorConfig.Up:
		if !RequestsAbove(e) && !e.Requests[e.Floor][elevatorConfig.HallUp] {
			e.Requests[e.Floor][elevatorConfig.HallDown] = false
			ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.HallDown}
			ElevatorRequestButtonLight(e.Floor, elevatorConfig.HallDown, false)
		}
		e.Requests[e.Floor][elevatorConfig.HallUp] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.HallUp}
		ElevatorRequestButtonLight(e.Floor, elevatorConfig.HallUp, false)

	case elevatorConfig.Down:
		if !RequestsBelow(e) && !e.Requests[e.Floor][elevatorConfig.HallDown] {
			e.Requests[e.Floor][elevatorConfig.HallUp] = false
			ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.HallUp}
			ElevatorRequestButtonLight(e.Floor, elevatorConfig.HallUp, false)
			
			
		}
		e.Requests[e.Floor][elevatorConfig.HallDown] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.HallDown}
		ElevatorRequestButtonLight(e.Floor, elevatorConfig.HallDown, false)

	case elevatorConfig.Stop:
		fallthrough
	default:
		e.Requests[e.Floor][elevatorConfig.HallUp] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.HallUp}
		ElevatorRequestButtonLight(e.Floor, elevatorConfig.HallUp, false)
		e.Requests[e.Floor][elevatorConfig.HallDown] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor : e.Floor, Button: elevatorConfig.HallDown}
		ElevatorRequestButtonLight(e.Floor, elevatorConfig.HallDown, false)
	}


	return e
}
