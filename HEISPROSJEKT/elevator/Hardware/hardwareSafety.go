package Hardware

import (
	"Driver-go/elevio"
	"fmt"
)

var _numFloors int = 4

func MotorDriection(motorDirection chan elevio.MotorDirection){
	for {
	d := <- motorDirection
	elevio.SetMotorDirection(d)
	}
}


func HardwareSafetyFeatures(pollObstructionChannel chan bool, pollStopButtonChannel chan bool, motorDirection chan elevio.MotorDirection){
 	for {
		select{

		case obstructionActivated:= <- pollObstructionChannel:

			fmt. Println(obstructionActivated)
				if obstructionActivated {

					motorDirection <- elevio.MD_Stop
					
				} else {
					motorDirection <- elevio.MD_Up
				}
		
			//when true 
			//case 1: elevator on a floor
				//keep door open until obstruction  is celared
				

			//case 2: elevator not on flor
				//ignore obstruction entirely
			
			//when false
				//case 1: obstruction was NOT recently activated
					//do nothing
				//case 2: obstruction WAS recently activated
					//keep door open for 3 more seconds, then normal operation



		case stopActivated := <- pollStopButtonChannel:
			fmt. Println(stopActivated)

			if (stopActivated){

				elevio.SetStopLamp(true)
				motorDirection <- elevio.MD_Stop
				TurnOffAllLights()

			} else {

				elevio.SetStopLamp(false)
				motorDirection <- elevio.MD_Up

			}
			


			//when true

			//when false
				//Case 1: Stop was NOT recently activated 
					//do nothing
				//case 2: Stop was recently activated
					//if on floor: Keep door open for 3 more seconds -> normal operation
					//else: continue normal operation

			}
	}
}