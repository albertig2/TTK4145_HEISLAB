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

func tunrOnOrderLightsBasedOnLocalQueue(elevator elevatorConfig.Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
			if elevator.LocalOrderQueue[floor][button] { //only set light if true, (clearing here will also mess withthe network order lights)
				orderButtonLight(floor, elevatorConfig.Button(button), elevator.LocalOrderQueue[floor][button])
			}
		}
	}
}

func turnOffAllOrderLights() {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := 0; button < elevatorConfig.N_BUTTONS; button++ {
			orderButtonLight(floor, elevatorConfig.Button(button), false)
		}
	}
}

func clearListOfOrderLighst(orderLightList []elevatorConfig.ButtonEvent) {
	for _, orderLightToBeCleard := range orderLightList {
		orderButtonLight(orderLightToBeCleard.Floor, orderLightToBeCleard.Button, false)
	}
}
