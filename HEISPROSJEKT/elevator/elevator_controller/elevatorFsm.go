package elevatorHardware

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func RunElevatorFsm(elevatorID string, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels, orderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	doorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	doorTimer.Stop()

	motorTimeoutTimer := time.NewTimer(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	motorTimeoutTimer.Stop()

	

	elevatorObject := InitializeElevator(elevatorID)



	InitElevatorHardware(&elevatorObject, motorTimeoutTimer)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:
			// Should wait for assignment before opening the door
			motorTimeoutTimer.Stop()
			motorTimeoutTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)

			HandleOnFloorArrival(&elevatorObject, doorTimer, floor, orderChannels.ServicedOrderChannel, motorTimeoutTimer)

		case recievedOrder := <-hardwareChannels.PollOrderButtonsChannel:

			order := elevatorConfig.ButtonEvent{Floor: recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Println("New order from FSM", order)
			orderChannels.NewRecievedOrderChannel <- order
			fmt.Println("Passed neworderchannel")

		case assignedOrder := <-orderChannels.NewAssignedOrderChannel:

			HandleRequestButtonPress(&elevatorObject, doorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button), orderChannels.ServicedOrderChannel, motorTimeoutTimer)

		case assignedPeerOrder := <-orderChannels.NewAssignedPeerOrderChannel:
			HandlelightSettingForPeerOrders(assignedPeerOrder.Floor, assignedPeerOrder.Button, true)
			//HandlePeerAssignedOrder(&elevatorObject, doorTimer, int(assignedPeerOrder.Floor), elevatorConfig.Button(assignedPeerOrder.Button), orderChannels.ServicedOrderChannel, motorTimeoutTimer)
			fmt.Printf("Turned on the light for %v at floor %v \n", elevatorConfig.ButtonToString(assignedPeerOrder.Button), assignedPeerOrder.Floor)
		case servicedPeerOrder := <-orderChannels.ServicedPeerOrderChannel:
			HandlelightSettingForPeerOrders(servicedPeerOrder.Floor, servicedPeerOrder.Button, false)
			fmt.Printf("Turned of the light for %v at floor %v \n", elevatorConfig.ButtonToString(servicedPeerOrder.Button), servicedPeerOrder.Floor)
			//HandlePeerServicedOrder(&elevatorObject, doorTimer, int(servicedPeerOrder.Floor), elevatorConfig.Button(servicedPeerOrder.Button), orderChannels.ServicedOrderChannel, motorTimeoutTimer)

		case stopActivated := <-hardwareChannels.PollStopButtonChannel:

			HandleStopButtonActivated(stopActivated, &elevatorObject, doorTimer, orderChannels.ServicedOrderChannel, motorTimeoutTimer)

		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:

			HandleObstructionActivated(obstructionActivated, &elevatorObject, doorTimer, orderChannels.ServicedOrderChannel, motorTimeoutTimer, hardwareChannels, synchronisationChannels)

		case <-doorTimer.C:
			OnDoorTimeout(&elevatorObject, doorTimer, orderChannels.ServicedOrderChannel, motorTimeoutTimer)

		case <-motorTimeoutTimer.C:
			HandleMotorTimeout(&elevatorObject, motorTimeoutTimer, hardwareChannels, synchronisationChannels)
		}

		select {
		case hardwareChannels.ElevatorObjectChannel <- elevatorObject:

		default:
		}
	}
}
