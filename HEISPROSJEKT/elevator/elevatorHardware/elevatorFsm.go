package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

//testcomment
func RunElevatorFsm(elevatorID string, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, orderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	doorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	doorTimer.Stop()

	elevatorObject := InitializeElevator(elevatorID)

	InitElevatorHardware(&elevatorObject)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:


			HandleOnFloorArrival(&elevatorObject, doorTimer, floor, hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)

		case recievedOrder := <-hardwareChannels.PollOrderButtonsChannel:

			order:= elevatorConfig.ButtonEvent{Floor : recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Println("New order from FSM", order)
			orderChannels.NewRecievedOrderChannel <- order


		case assignedOrder := <- orderChannels.NewAssignedOrderChannel:

			HandleRequestButtonPress(&elevatorObject, doorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button), hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)

		case stopActivated := <-hardwareChannels.PollStopButtonChannel:

			HandleStopButtonActivated( stopActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)

		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
	
			HandleObstructionActivated(obstructionActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)



		case <-doorTimer.C:
			OnDoorTimeout(&elevatorObject, doorTimer, orderChannels.ServicedOrderChannel)

		}

		select {
		case hardwareChannels.ElevatorObjectChannel <- elevatorObject:
		default:
		}
	}
}
