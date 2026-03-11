package elevatorHardware

import (
	//"Driver-go/elevio"
	//"HEISPROSJEKT/debuggingHelpers"
	//"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	//"fmt"
	"time"
)

//testcomment
func RunElevatorFsm(elevatorID string, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, orderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
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

			HandleOnFloorArrival(&elevatorObject, doorTimer, floor, hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)

		case recievedOrder := <-hardwareChannels.PollOrderButtonsChannel:
			//orderChannels.NewRecievedOrderChannel <- floor, buttontype (butten event)
			order:= elevatorConfig.ButtonEvent{Floor : recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Println("New order from FSM", order)
			orderChannels.NewRecievedOrderChannel <- order

			//HandleRequestButtonPress(&elevatorObject, doorTimer, int(recievedOrder.Floor), elevatorConfig.Button(recievedOrder.Button), hardwareChannels.MotorDirectionChannel)
		case assignedOrder := <- orderChannels.NewAssignedOrderChannel:

			HandleRequestButtonPress(&elevatorObject, doorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button), hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)

		case stopActivated := <-hardwareChannels.PollStopButtonChannel:

			HandleStopButtonActivated( stopActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)

		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
	
			HandleObstructionActivated(obstructionActivated, &elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel,orderChannels.ServicedOrderChannel)


		// orderdque <- new order quqe (REq_matrix)
		case <-doorTimer.C:
			OnDoorTimeout(&elevatorObject, doorTimer, orderChannels.ServicedOrderChannel)

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
