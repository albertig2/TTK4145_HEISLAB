package Hardware

import (
	"Driver-go/elevio"
	//"time"
	// "fmt"
)

func TurnOffAllOrderLights() {
	for floor := 0; floor < _numFloors; floor++ {
		for button := elevio.ButtonType(0); button < 3; button++ {
			elevio.SetButtonLamp(button, floor, false)
		}
	}
}
