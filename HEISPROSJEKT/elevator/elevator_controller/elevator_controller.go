package elevatorController

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func LocalElevatorController(ownId string, controllerChannels elevatorConfig.ElevatorControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels, orderChannels elevatorConfig.OrderChannels) {
	openDoorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	openDoorTimer.Stop()

	detectMotorFailureTimer := time.NewTimer(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	detectMotorFailureTimer.Stop()

	elevator := initializeEmptyElevator(ownId)

	initializeElevatorHardware(&elevator, detectMotorFailureTimer)

	for {
		select {
		case floor := <-controllerChannels.FloorSensorChannel:
			// Should wait for assignment before opening the door
			detectMotorFailureTimer.Stop()
			detectMotorFailureTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)

			handleOnFloorArrival(&elevator, openDoorTimer, floor, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case recievedOrder := <-controllerChannels.PollOrderButtonsChannel:
			orderChannels.NewRecievedOrderChannel <- elevatorConfig.ButtonEvent{Floor: recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Printf("New order from FSM: (%v , %v) \n", elevatorConfig.ButtonToString(elevatorConfig.Button(int(recievedOrder.Button))), recievedOrder.Floor)

		case assignedOrder := <-orderChannels.NewAssignedOrderChannel:
			handleRequestButtonPressd(&elevator, openDoorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button), orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case assignedPeerOrder := <-orderChannels.NewAssignedPeerOrderChannel:
			handleLightSettingForPeerOrders(assignedPeerOrder.Floor, assignedPeerOrder.Button, true)
			fmt.Printf("Turned on the light for %v at floor %v \n", elevatorConfig.ButtonToString(assignedPeerOrder.Button), assignedPeerOrder.Floor)

		case servicedPeerOrder := <-orderChannels.ServicedPeerOrderChannel:
			handleLightSettingForPeerOrders(servicedPeerOrder.Floor, servicedPeerOrder.Button, false)
			fmt.Printf("Turned of the light for %v at floor %v \n", elevatorConfig.ButtonToString(servicedPeerOrder.Button), servicedPeerOrder.Floor)

		case stopActivated := <-controllerChannels.PollStopButtonChannel:
			handleStopButton(stopActivated, &elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case obstructionActivated := <-controllerChannels.PollObstructionChannel:
			handleObstruction(obstructionActivated, &elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer, controllerChannels, synchronisationChannels)

		case <-openDoorTimer.C:
			handleDoorTimeout(&elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case <-detectMotorFailureTimer.C:
			handleDetectedMotorFailure(&elevator, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
		}
		select {
		case controllerChannels.LocalElevatorChannel <- elevator:
		default:
		}
	}
}
