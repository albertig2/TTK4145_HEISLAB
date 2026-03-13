package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronisation"
	//"fmt"
	"time"
)

func RunElevatorFsm(elevatorID string, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels synchronisation.SynchronisationChannels) {
	doorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	doorTimer.Stop()

	motorTimeoutTimer :=  time.NewTimer(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	motorTimeoutTimer.Stop()

	elevatorObject := InitializeElevator(elevatorID)

	InitElevatorHardware(&elevatorObject, motorTimeoutTimer)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:

			HandleOnFloorArrival(&elevatorObject, doorTimer, floor, hardwareChannels.MotorDirectionChannel, motorTimeoutTimer)

		case recievedOrder := <-hardwareChannels.PollOrderButtonsChannel:

			order:= elevatorConfig.ButtonEvent{Floor : recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Println("New order from FSM", order)
			orderChannels.NewRecievedOrderChannel <- order


		case assignedOrder := <- orderChannels.NewAssignedOrderChannel:

			HandleRequestButtonPress(&elevatorObject, doorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button),orderChannels.ServicedOrderChannel)

		case stopActivated := <-hardwareChannels.PollStopButtonChannel:

			HandleStopButtonActivated( stopActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel, motorTimeoutTimer)

		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
	
			HandleObstructionActivated(obstructionActivated, &elevatorObject, doorTimer, orderChannels.ServicedOrderChannel)

		case <-doorTimer.C:
			OnDoorTimeout(&elevatorObject, doorTimer, motorTimeoutTimer)

		

		case <-motorTimeoutTimer.C:
			HandleMotorFailure(&elevatorObject, motorTimeoutTimer, hardwareChannels, synchronisationChannels)
			
		}

		
		select {
		case hardwareChannels.ElevatorObjectChannel <- elevatorObject:
		
		default:
		}
	}
}
