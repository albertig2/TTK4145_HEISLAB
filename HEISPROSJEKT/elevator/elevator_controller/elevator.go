package elevatorController

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"time"
)

//------------------Sensors-------------------------------

func ElevatorFloorSensor() int {
	return elevio.GetFloor()
}

//-------------------motor------------------------------------


func ElevatorMotorDirection(d elevatorConfig.Direction, motorTimeoutTimer *time.Timer) {
	elevio.SetMotorDirection(elevio.MotorDirection(d))

	if d == 0 {
		motorTimeoutTimer.Stop()
	} else {
		motorTimeoutTimer.Stop()
		motorTimeoutTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	}
}


// ------------------Buttons--------------------------------
func ElevatorRequestButton(f int, b elevatorConfig.Button) bool {
	return elevio.GetButton((elevio.ButtonType)(b), f)
}

func ElevatorStopButton() bool {
	return elevio.GetStop()
}

func ElevatorObstruction() bool {
	return elevio.GetObstruction()
}

//-------------------Lights----------------------------------

func ElevatorFloorIndicatorLight(f int) {
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

func SetAllLights(elevator elevatorConfig.Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if (elevator.Requests[floor][btn]){ //only set light if true, (clearing here will also mess withthe network order lights)

			ElevatorRequestButtonLight(floor, elevatorConfig.Button(btn), elevator.Requests[floor][btn])
			}
		}
	}
}
func TurnOffAllOrderLights() {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := elevio.ButtonType(0); button < 3; button++ {
			elevio.SetButtonLamp(button, floor, false)
		}
	}
}



