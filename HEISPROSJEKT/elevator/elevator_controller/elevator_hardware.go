package elevatorController

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"time"
)

/*
The elevator hardware sub module functions as an intermediate layer between the higher-level 
code and the hardware. It is essentially made to “hide” the raw hardware layer of the elevio 
hardware driver (elevio is pre made, and handed out as a project resource), and support the 
error handling in the finite state machine.
*/


//------------------Sensors-------------------------------

func FloorSensor() int {
	return elevio.GetFloor()
}

//-------------------motor------------------------------------

func motorDirection(direction elevatorConfig.Direction, detectMotorFailureTimer *time.Timer) {
	elevio.SetMotorDirection(elevio.MotorDirection(direction))

	if direction == 0 {
		detectMotorFailureTimer.Stop()
	} else {
		detectMotorFailureTimer.Stop()
		detectMotorFailureTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	}
}

// ------------------Buttons--------------------------------
func orderButton(floor int, button elevatorConfig.Button) bool {
	return elevio.GetButton((elevio.ButtonType)(button), floor)
}

func stopButton() bool {
	return elevio.GetStop()
}

func obstruction() bool {
	return elevio.GetObstruction()
}

//-------------------Lights----------------------------------

func floorIndicatorLight(floor int) {
	elevio.SetFloorIndicator(floor)
}

func orderButtonLight(floor int, button elevatorConfig.Button, lightValue bool) {
	elevio.SetButtonLamp(elevio.ButtonType(button), floor, lightValue)
}

func doorLight(lightValue bool) {
	elevio.SetDoorOpenLamp(lightValue)

}

func stopButtonLight(lightValue bool) {
	elevio.SetStopLamp(lightValue)
}

func setAllOrderLights(elevator elevatorConfig.Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if elevator.LocalOrderQueue[floor][btn] { //only set light if true, (clearing here will also mess withthe network order lights)

				orderButtonLight(floor, elevatorConfig.Button(btn), elevator.LocalOrderQueue[floor][btn])
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
