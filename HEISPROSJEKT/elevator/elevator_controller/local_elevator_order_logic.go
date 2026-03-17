package elevatorController

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
			if elevator.LocalOrderQueue[floor][button] {
				return true
			}
		}
	}
	return false
}

func ordersBelowCurrentFloor(elevator elevatorConfig.Elevator) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
			if elevator.LocalOrderQueue[floor][button] {
				return true
			}
		}
	}
	return false
}

func ordersAtCurrentFloor(elevator elevatorConfig.Elevator) bool {
	for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
		if elevator.LocalOrderQueue[elevator.Floor][button] {
			return true
		}
	}
	return false
}

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
		return elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallDown] ||
			elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.Cab] ||
			!ordersBelowCurrentFloor(elevator)

	case elevatorConfig.Up:
		return elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallUp] ||
			elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.Cab] ||
			!ordersAboveCurrentFloor(elevator)

	case elevatorConfig.Stop:
		fallthrough
	default:
		return true
	}
}

func shouldClearOrderImmediately(elevator elevatorConfig.Elevator, buttonFloor int, buttonType elevatorConfig.Button) bool {
	return elevator.Floor == buttonFloor &&
		((elevator.Direction == elevatorConfig.Up && buttonType == elevatorConfig.HallUp) ||
			(elevator.Direction == elevatorConfig.Down && buttonType == elevatorConfig.HallDown) ||
			elevator.Direction == elevatorConfig.Stop ||
			buttonType == elevatorConfig.Cab)
}

func clearOrdersAtCurrentFloor(elevator elevatorConfig.Elevator, ServicedOrderChannel chan elevatorConfig.ButtonEvent) (elevatorConfig.Elevator, []elevatorConfig.ButtonEvent) {

	clearedOrders := []elevatorConfig.ButtonEvent {}

	elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.Cab] = false
	clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.Cab})

	// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.Cab} //clear light after the order after network is pinged
	// orderButtonLight(elevator.Floor, elevatorConfig.Cab, false)

	switch elevator.Direction {
	case elevatorConfig.Up:
		if !ordersAboveCurrentFloor(elevator) && !elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallUp] {
			elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallDown] = false
			clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown})

			// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown}
			// orderButtonLight(elevator.Floor, elevatorConfig.HallDown, false)
		}

		elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallUp] = false
		clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp})

		// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp}
		// orderButtonLight(elevator.Floor, elevatorConfig.HallUp, false)

	case elevatorConfig.Down:
		if !ordersBelowCurrentFloor(elevator) && !elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallDown] {
			
			elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallUp] = false
			clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp})
			
			// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp}
			// orderButtonLight(elevator.Floor, elevatorConfig.HallUp, false)

		}
		elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallDown] = false
		clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown})

		// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown}
		// orderButtonLight(elevator.Floor, elevatorConfig.HallDown, false)

	case elevatorConfig.Stop:
		fallthrough
	default:
		elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallUp] = false
		clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp})

		// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallUp}
		// orderButtonLight(elevator.Floor, elevatorConfig.HallUp, false)

		elevator.LocalOrderQueue[elevator.Floor][elevatorConfig.HallDown] = false
		clearedOrders = append(clearedOrders, elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown})

		// ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: elevator.Floor, Button: elevatorConfig.HallDown}
		// orderButtonLight(elevator.Floor, elevatorConfig.HallDown, false)
	}

	return elevator, clearedOrders
}
