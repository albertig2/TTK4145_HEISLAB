package elevatorHardware

import (
	//"Driver-go/elevio"
	//"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	//"fmt"
	"time"
)
//testcomment
func RunElevatorFsm(elevatorID string, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt) {
	doorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	doorTimer.Stop()

	//motorTimer =  time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	//motortimer.Stop()
	elevatorObject := InitializeElevator(elevatorID)

	InitElevatorHardware(&elevatorObject)

	//go updateMotorDirection(hardwareChannels.MotorDirectionChannel)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:
			//serviced order <- floor, order (elevatorsystem)

			HandleOnFloorArrival(&elevatorObject, doorTimer, floor, hardwareChannels.MotorDirectionChannel)

		case recievedOrder := <-hardwareChannels.PollOrderButtonsChannel:
			//neworder <- floor, buttontype (butten event)

			HandleRequestButtonPress(&elevatorObject, doorTimer, int(recievedOrder.Floor), elevatorConfig.Button(recievedOrder.Button), hardwareChannels.MotorDirectionChannel)

		case stopActivated := <-hardwareChannels.PollStopButtonChannel:

			HandleStopButtonActivated( stopActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel)

		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
	
			HandleObstructionActivated(obstructionActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel)


		// orderdque <- new order quqe (REq_matrix)
		case <-doorTimer.C:
			OnDoorTimeout(&elevatorObject, doorTimer)

		}

		// case <-motortimer.C:
		// 	HandleMotorFailur(&elevatorObject, doorTimer)

		// }

		
		select {
		case hardwareChannels.ElevatorObjectChannel <- elevatorObject:
		default:
		}
	}
}
