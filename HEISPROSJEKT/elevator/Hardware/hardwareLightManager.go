package Hardware

import (
	"Driver-go/elevio"
	// "fmt"
)


func TurnOffAllLights(){

	for floor := 0; floor < _numFloors; floor++ {
		for button := elevio.ButtonType(0); button< 3; button++ {
			elevio.SetButtonLamp(button, floor, false)
		}

	}
}