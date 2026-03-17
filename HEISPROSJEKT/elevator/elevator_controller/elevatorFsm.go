package elevatorController

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func RunElevatorFsm(ownId string, controllerChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels, orderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	openDoorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	openDoorTimer.Stop()

	detectMotorFailureTimer := time.NewTimer(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	detectMotorFailureTimer.Stop()

	elevatorObject := initializeElevator(ownId)

	initializeElevatorHardware(&elevatorObject, detectMotorFailureTimer)

	for {
		select {
		case floor := <-controllerChannels.FloorSensorChannel:
			// Should wait for assignment before opening the door
			detectMotorFailureTimer.Stop()
			detectMotorFailureTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)

			handleOnFloorArrival(&elevatorObject, openDoorTimer, floor, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case recievedOrder := <-controllerChannels.PollOrderButtonsChannel:

			orderChannels.NewRecievedOrderChannel <- elevatorConfig.ButtonEvent{Floor: recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Printf("New order from FSM: ( %v , %v)", elevatorConfig.ButtonToString(elevatorConfig.Button(int(recievedOrder.Button))), recievedOrder.Floor)

		case assignedOrder := <-orderChannels.NewAssignedOrderChannel:

			handleRequestButtonPressd(&elevatorObject, openDoorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button), orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case assignedPeerOrder := <-orderChannels.NewAssignedPeerOrderChannel:

			handleLightSettingForPeerOrders(assignedPeerOrder.Floor, assignedPeerOrder.Button, true)

			fmt.Printf("Turned on the light for %v at floor %v \n", elevatorConfig.ButtonToString(assignedPeerOrder.Button), assignedPeerOrder.Floor)
		case servicedPeerOrder := <-orderChannels.ServicedPeerOrderChannel:

			handleLightSettingForPeerOrders(servicedPeerOrder.Floor, servicedPeerOrder.Button, false)

			fmt.Printf("Turned of the light for %v at floor %v \n", elevatorConfig.ButtonToString(servicedPeerOrder.Button), servicedPeerOrder.Floor)
		case stopActivated := <-controllerChannels.PollStopButtonChannel:

			handleStopButton(stopActivated, &elevatorObject, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case obstructionActivated := <-controllerChannels.PollObstructionChannel:

			handleObstruction(obstructionActivated, &elevatorObject, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer, controllerChannels, synchronisationChannels)

		case <-openDoorTimer.C:
			handleDoorTimeout(&elevatorObject, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case <-detectMotorFailureTimer.C:
			handleDetectedMotorFailure(&elevatorObject, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
		}
		select {
		case controllerChannels.ElevatorObjectChannel <- elevatorObject:

		default:
		}
	}
}
