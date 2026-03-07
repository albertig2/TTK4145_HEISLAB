package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
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
	requests  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool
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

func Elevator_print(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevatorConfig.DirectionToString(es.direction))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", Elevator_behaviorToString(es.behaviour))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevatorConfig.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if (f == elevatorConfig.N_FLOORS-1 && elevatorConfig.Button(btn) == elevatorConfig.HallUp) || (f == 0 && elevatorConfig.Button(btn) == elevatorConfig.HallDown) {
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
	elevio.Init("localhost:15657", elevatorConfig.N_FLOORS)
	es := Elevator{floor: -1, direction: elevatorConfig.Stop, behaviour: EB_Idle, config: Config{doorOpenDuration_s: 3.0}}
	return es
}

func Elevator_floorSensor() int {
	return elevio.GetFloor()
}

func Elevator_requestButton(f int, b elevatorConfig.Button) bool {
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

func Elevator_requestButtonLight(f int, b elevatorConfig.Button, v bool) {
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
