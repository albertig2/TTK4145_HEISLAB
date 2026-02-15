package main

import (
	"fmt"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

type Dirn int

const (
	D_Down Dirn = -1
	D_Stop Dirn = 0
	D_Up   Dirn = 1
)

type Button int

const (
	B_HallUp   Button = 0
	B_HallDown Button = 1
	B_Cab      Button = 2
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2
)

type Config struct {
	doorOpenDuration_s float64
}

type Elevator struct {
	floor     int
	dirn      Dirn
	requests  [N_FLOORS][N_BUTTONS]bool
	behaviour ElevatorBehaviour
	config    Config
}

func elevator_behaviorToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

func elevator_dirnToString(d Dirn) string {
	switch d {
	case D_Up:
		return "D_Up"
	case D_Down:
		return "D_Down"
	case D_Stop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}

func elevator_buttonToString(b Button) string {
	switch b {
	case B_HallUp:
		return "B_HallUp"
	case B_HallDown:
		return "B_HallDown"
	case B_Cab:
		return "B_Cab"
	default:
		return "B_UNDEFINED"
	}
}

func elevator_print(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevator_dirnToString(es.dirn))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", elevator_behaviorToString(es.behaviour))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < N_BUTTONS; btn++ {
			if (f == N_FLOORS-1 && Button(btn) == B_HallUp) || (f == 0 && Button(btn) == B_HallDown) {
				fmt.Printf("|     ")
			} else {
				if es.requests[f][btn] {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func elevator_uninitialized() Elevator {
	elevator_hardware_init()
	es := Elevator{floor: -1, dirn: D_Stop, behaviour: EB_Idle, config: Config{doorOpenDuration_s: 3.0}}
	return es
}

func elevator_floorSensor() int {
	return elevator_hardware_get_floor_sensor_signal()
}

func elevator_requestButton(f int, b Button) int {
	return elevator_hardware_get_button_signal((elevator_hardware_button_type_t)(b), f)
}

func elevator_stopButton() int {
	return elevator_hardware_get_stop_signal()
}

func elevator_obstruction() int {
	return elevator_hardware_get_obstruction_signal()
}

func elevator_floorIndicator(f int) {
	elevator_hardware_set_floor_indicator(f)
}

func elevator_requestButtonLight(f int, b Button, v int) {
	elevator_hardware_set_button_lamp((elevator_hardware_button_type_t)(b), f, v)
}

func elevator_doorLight(v int) {
	elevator_hardware_set_door_open_lamp(v)
}

func elevator_stopButtonLight(v int) {
	elevator_hardware_set_stop_lamp(v)
}

func elevator_motorDirection(d Dirn) {
	elevator_hardware_set_motor_direction((elevator_hardware_motor_direction_t)(d))
}
