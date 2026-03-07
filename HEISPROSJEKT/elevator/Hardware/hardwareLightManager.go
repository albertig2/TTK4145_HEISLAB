package hardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	//"time"
	// "fmt"
)

func TurnOffAllOrderLights() {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := elevio.ButtonType(0); button < 3; button++ {
			elevio.SetButtonLamp(button, floor, false)
		}
	}
}
