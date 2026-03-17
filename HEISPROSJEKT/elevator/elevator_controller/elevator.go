package elevatorController

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"time"
)

//------------------Sensors-------------------------------

func FloorSensor() int {
	return elevio.GetFloor()
}

//-------------------motor------------------------------------

func motorDirection(d elevatorConfig.Direction, detectMotorFailureTimer *time.Timer) {
	elevio.SetMotorDirection(elevio.MotorDirection(d))

	if d == 0 {
		detectMotorFailureTimer.Stop()
	} else {
		detectMotorFailureTimer.Stop()
		detectMotorFailureTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	}
}

// ------------------Buttons--------------------------------
func orderButton(f int, b elevatorConfig.Button) bool {
	return elevio.GetButton((elevio.ButtonType)(b), f)
}

func stopButton() bool {
	return elevio.GetStop()
}

func obstruction() bool {
	return elevio.GetObstruction()
}

//-------------------Lights----------------------------------

func floorIndicatorLight(f int) {
	elevio.SetFloorIndicator(f)
}

func orderButtonLight(f int, b elevatorConfig.Button, v bool) {
	elevio.SetButtonLamp(elevio.ButtonType(b), f, v)
}

func doorLight(v bool) {
	elevio.SetDoorOpenLamp(v)

}

func stopButtonLight(v bool) {
	elevio.SetStopLamp(v)
}

func setAllOrderLights(elevator elevatorConfig.Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if elevator.Requests[floor][btn] { //only set light if true, (clearing here will also mess withthe network order lights)

				orderButtonLight(floor, elevatorConfig.Button(btn), elevator.Requests[floor][btn])
			}
		}
	}
}
func turnOffAllOrderLights() {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := elevio.ButtonType(0); button < 3; button++ {
			elevio.SetButtonLamp(button, floor, false)
		}
	}
}
