package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
)

type directionBehaviorPair struct {
	direction elevatorConfig.Direction
	behavior  elevatorConfig.Behavior
}

func ordersAboveCurrentFloor(elevator elevatorConfig.Elevator) bool {
	for floor := elevator.Floor + 1; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
			if elevator.Requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func ordersBelowCurrentFloor(elevator elevatorConfig.Elevator) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
			if elevator.Requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func ordersAtCurrentFloor(elevator elevatorConfig.Elevator) bool {
	for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
		if elevator.Requests[elevator.Floor][button] {
			return true
		}
	}
	return false
}

// “API” lik headeren: behold navnet
func chooseDirectionBasedOnOrders(elevator elevatorConfig.Elevator) directionBehaviorPair {
	switch elevator.Direction {
	case elevatorConfig.Up:
		if ordersAboveCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		if ordersAtCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Down, elevatorConfig.DoorOpen}
		}
		if ordersBelowCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		return directionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	case elevatorConfig.Down:
		if ordersBelowCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		if ordersAtCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Up, elevatorConfig.DoorOpen}
		}
		if ordersAboveCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		return directionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	case elevatorConfig.Stop:
		if ordersAtCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Stop, elevatorConfig.DoorOpen}
		}
		if ordersAboveCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Up, elevatorConfig.Moving}
		}
		if ordersBelowCurrentFloor(elevator) {
			return directionBehaviorPair{elevatorConfig.Down, elevatorConfig.Moving}
		}
		return directionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}

	default:
		return directionBehaviorPair{elevatorConfig.Stop, elevatorConfig.Idle}
	}
}

func shouldStopAtCurrentFloor(elevator elevatorConfig.Elevator) bool {
	switch elevator.Direction {
	case elevatorConfig.Down:
		return elevator.Requests[elevator.Floor][elevatorConfig.HallDown] ||
			elevator.Requests[elevator.Floor][elevatorConfig.Cab] ||
			!ordersBelowCurrentFloor(elevator)

	case elevatorConfig.Up:
		return elevator.Requests[elevator.Floor][elevatorConfig.HallUp] ||
			elevator.Requests[elevator.Floor][elevatorConfig.Cab] ||
			!ordersAboveCurrentFloor(elevator)

	case elevatorConfig.Stop:
		fallthrough
	default:
		return true
	}
}

func shouldClearOrderImmediately(elevator elevatorConfig.Elevator, button_Floor int, button_type elevatorConfig.Button) bool {
	return elevator.Floor == button_Floor &&
		((elevator.Direction == elevatorConfig.Up && button_type == elevatorConfig.HallUp) ||
			(elevator.Direction == elevatorConfig.Down && button_type == elevatorConfig.HallDown) ||
			elevator.Direction == elevatorConfig.Stop ||
			button_type == elevatorConfig.Cab)
}

func clearOrdersAtCurrentFloor(elevator elevatorConfig.Elevator, ServicedOrderChannel chan elevatorConfig.ButtonEvent) elevatorConfig.Elevator {

	elevator.Requests[elevator.Floor][elevatorConfig.Cab] = false
	ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.Cab} //clear light after the order after network is pinged
	ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.Cab, false)

	switch elevator.Direction {
	case elevatorConfig.Up:
		if !ordersAboveCurrentFloor(elevator) && !elevator.Requests[elevator.Floor][elevatorConfig.HallUp] {
			elevator.Requests[elevator.Floor][elevatorConfig.HallDown] = false
			ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown}
			ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.HallDown, false)
		}
		elevator.Requests[elevator.Floor][elevatorConfig.HallUp] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp}
		ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.HallUp, false)

	case elevatorConfig.Down:
		if !ordersBelowCurrentFloor(elevator) && !elevator.Requests[elevator.Floor][elevatorConfig.HallDown] {
			elevator.Requests[elevator.Floor][elevatorConfig.HallUp] = false
			ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp}
			ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.HallUp, false)

		}
		elevator.Requests[elevator.Floor][elevatorConfig.HallDown] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown}
		ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.HallDown, false)

	case elevatorConfig.Stop:
		fallthrough
	default:
		elevator.Requests[elevator.Floor][elevatorConfig.HallUp] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp}
		ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.HallUp, false)
		elevator.Requests[elevator.Floor][elevatorConfig.HallDown] = false
		ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown}
		ElevatorRequestButtonLight(elevator.Floor, elevatorConfig.HallDown, false)
	}

	return elevator
}
