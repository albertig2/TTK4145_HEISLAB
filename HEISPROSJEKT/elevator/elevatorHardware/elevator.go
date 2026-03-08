package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

type Config struct {
	DoorOpenDuration_s time.Duration
}

type Elevator struct {
	OwnId     string
	Floor     int
	Direction elevatorConfig.Direction
	Requests  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool
	Behavior  elevatorConfig.Behavior
	Config    Config
}

func InitializeElevatorObject(id string) Elevator {

	config := Config{
		DoorOpenDuration_s: 3 * time.Second,
	}

	elevator := Elevator{
		OwnId:     id,
		Floor:     elevio.GetFloor(),
		Direction:      elevatorConfig.Stop,
		Requests:  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
		Behavior: elevatorConfig.Idle,
		Config:    config,
	}
	return elevator
}

func Elevator_print(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.Floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevatorConfig.DirectionToString(es.Direction))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", elevatorConfig.BehaviorToString(es.Behavior))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevatorConfig.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if (f == elevatorConfig.N_FLOORS-1 && elevatorConfig.Button(btn) == elevatorConfig.HallUp) || (f == 0 && elevatorConfig.Button(btn) == elevatorConfig.HallDown) {
				fmt.Printf("|     ")
			} else {
				if es.Requests[f][btn] {
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
	es := Elevator{Floor: -1, Direction: elevatorConfig.Stop, Behavior: elevatorConfig.Idle, Config: Config{DoorOpenDuration_s: 3.0}}
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
