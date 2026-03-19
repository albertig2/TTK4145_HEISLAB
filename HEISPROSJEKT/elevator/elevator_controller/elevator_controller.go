package elevatorController

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func LocalElevatorController(ownId string, controllerChannels elevatorConfig.ControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels, orderChannels elevatorConfig.OrderChannels) {
	openDoorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	openDoorTimer.Stop()

	detectMotorFailureTimer := time.NewTimer(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)
	detectMotorFailureTimer.Stop()

	sendUpdateToPeerViewTicker := time.NewTicker(time.Second / 10)

	elevator := initializeEmptyElevator(ownId)

	initializeElevatorHardware(&elevator, detectMotorFailureTimer)

	for {
		select {
		case floor := <-controllerChannels.FloorSensorChannel:
			detectMotorFailureTimer.Stop()
			detectMotorFailureTimer.Reset(elevatorConfig.MOTOR_TIMEOUT_DURATION_S)

			handleOnFloorArrival(&elevator, openDoorTimer, floor, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case recievedOrder := <-controllerChannels.PollOrderButtonsChannel:
			orderChannels.NewRecievedOrderChannel <- elevatorConfig.ButtonEvent{Floor: recievedOrder.Floor, Button: elevatorConfig.Button(recievedOrder.Button)}
			fmt.Printf("New order from FSM: (%v , %v) \n", elevatorConfig.ButtonToString(elevatorConfig.Button(int(recievedOrder.Button))), recievedOrder.Floor)

		case assignedOrder := <-orderChannels.NewAssignedOrderChannel:
			handleOrderButtonPressd(&elevator, openDoorTimer, int(assignedOrder.Floor), elevatorConfig.Button(assignedOrder.Button), orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case assignedPeerOrder := <-orderChannels.NewAssignedPeerOrderChannel:
			handleLightSettingForPeerOrders(assignedPeerOrder.Floor, assignedPeerOrder.Button, true)

		case servicedPeerOrder := <-orderChannels.ServicedPeerOrderChannel:
			handleLightSettingForPeerOrders(servicedPeerOrder.Floor, servicedPeerOrder.Button, false)

		case stopActivated := <-controllerChannels.PollStopButtonChannel:
			handleStopButton(stopActivated, &elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case obstructionActivated := <-controllerChannels.PollObstructionChannel:
			handleObstruction(obstructionActivated, &elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer, controllerChannels, synchronisationChannels)

		case <-openDoorTimer.C:
			handleDoorTimeout(&elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case <-detectMotorFailureTimer.C:
			handleDetectedMotorFailure(&elevator, detectMotorFailureTimer, controllerChannels, synchronisationChannels)

		case <-sendUpdateToPeerViewTicker.C:
			controllerChannels.LocalElevatorChannel <- elevator
		}
	}
}
