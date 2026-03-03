package hardwareSafety

import (
	"Driver-go/elevio"
	"fmt"
)



func hardwareSafetyFeatures(pollObstructionChannel chan bool, pollStopButtonChannel chan bool, motorDirection elevio.MotorDirection){

	select{
	case obstructionActivated:= <- pollObstructionChannel:
		fmt. Println(obstructionActivated)
			if obstructionActivated {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(motorDirection)
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
		for floor := 0; floor < numFloors; f++ {
			for button := elevio.ButtonType(0); b < 3; b++ {
				elevio.SetButtonLamp(buton, floor, false)
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

