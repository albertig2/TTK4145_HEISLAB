package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

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
	doorOpenDuration_s time.Duration
}

type Elevator struct {
	floor     int
	direction elevatorConfig.Direction
	requests  [N_FLOORS][N_BUTTONS]bool
	behaviour ElevatorBehaviour
	config    Config
}

func Elevator_behaviorToString(eb ElevatorBehaviour) string {
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

func Elevator_buttonToString(b Button) string {
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

func Elevator_print(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "direction", ElevatorDirectionToString(es.direction))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", Elevator_behaviorToString(es.behaviour))
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

func Elevator_uninitialized() Elevator {
	elevio.Init("localhost:15657", N_FLOORS)
	es := Elevator{floor: -1, direction: elevatorConfig.Stop, behaviour: EB_Idle, config: Config{doorOpenDuration_s: 3.0}}
	return es
}

func Elevator_floorSensor() int {
	return elevio.GetFloor()
}

func Elevator_requestButton(f int, b Button) bool {
	return elevio.GetButton((elevio.ButtonType)(b), f)
}

func Elevator_stopButton() bool {
	return elevio.GetStop()
}

func Elevator_obstruction() bool {
	return elevio.GetObstruction()
}

func Elevator_floorIndicator(f int) {
	elevio.SetFloorIndicator(f)
}

func Elevator_requestButtonLight(f int, b Button, v bool) {
	elevio.SetButtonLamp(elevio.ButtonType(b), f, v)
}

func Elevator_doorLight(v bool) {
	elevio.SetDoorOpenLamp(v)

}

func Elevator_stopButtonLight(v bool) {
	elevio.SetStopLamp(v)

}

func Elevator_motorDirection(d elevatorConfig.Direction) {
	elevio.SetMotorDirection(elevio.MotorDirection(d))

}
