package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"time"
)

func InitializeElevatorObject(id string) elevatorConfig.Elevator {

	config := elevatorConfig.Config{
		DoorOpenDuration_s: 3 * time.Second,
	}

	elevator := elevatorConfig.Elevator{
		OwnId:     id,
		Floor:     elevio.GetFloor(),
		Direction: elevatorConfig.Stop,
		Requests:  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
		Behavior:  elevatorConfig.Idle,
		Config:    config,
	}
	return elevator
}



func Elevator_uninitialized() elevatorConfig.Elevator {
	elevio.Init("localhost:15657", elevatorConfig.N_FLOORS)
	es := elevatorConfig.Elevator{Floor: -1, Direction: elevatorConfig.Stop, Behavior: elevatorConfig.Idle, Config: elevatorConfig.Config{DoorOpenDuration_s: 3.0}}
	return es
}

func ElevatorFloorSensor() int {
	return elevio.GetFloor()
}

func ElevatorRequestButton(f int, b elevatorConfig.Button) bool {
	return elevio.GetButton((elevio.ButtonType)(b), f)
}

func ElevatorStopButton() bool {
	return elevio.GetStop()
}

func ElevatorObstruction() bool {
	return elevio.GetObstruction()
}

func ElevatorFloorIndicator(f int) {
	elevio.SetFloorIndicator(f)
}

func ElevatorRequestButtonLight(f int, b elevatorConfig.Button, v bool) {
	elevio.SetButtonLamp(elevio.ButtonType(b), f, v)
}

func ElevatorDoorLight(v bool) {
	elevio.SetDoorOpenLamp(v)

}

func ElevatorStopButtonLight(v bool) {
	elevio.SetStopLamp(v)

}

func ElevatorMotorDirection(d elevatorConfig.Direction) {
	elevio.SetMotorDirection(elevio.MotorDirection(d))
}
