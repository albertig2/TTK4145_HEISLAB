package elevatorController

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"fmt"
	"time"
)

/*
This is file contains one global function, LocalElevatorController.This is the finite state machine
loop for handleing the local elevator behavior. The controller loop contains 3 objects: One elevator, which contains the internal
state variables needed to make decisions in the state machine and two timers. One timer is used to detect
motor failure, and one is used to control the amount of time the door is open. The infinite for loop contains
two select cases. One is used to detect events from the hardware and the order assigner and trigger the correct
handler corresponding to the event. The other case sends the elevator object to the synchronization module.
*/

func LocalElevatorController(ownId string, controllerChannels elevatorConfig.ControllerChannels, synchronizationChannels elevatorConfig.SynchronizationChannels, orderChannels elevatorConfig.OrderChannels) {
	openDoorTimer := time.NewTimer(elevatorConfig.DoorOpenDurationInSeconds)
	openDoorTimer.Stop()

	detectMotorFailureTimer := time.NewTimer(elevatorConfig.MotorTimeOutDurationInSeconds)
	detectMotorFailureTimer.Stop()

	sendElevatorUpdateTicker := time.NewTicker(time.Second / 10)

	elevator := initializeEmptyElevator(ownId)

	initializeElevatorHardware(&elevator, detectMotorFailureTimer)

	for {
		select {
		case floor := <-controllerChannels.PollFloorSensorChannel:
			detectMotorFailureTimer.Stop()
			detectMotorFailureTimer.Reset(elevatorConfig.MotorTimeOutDurationInSeconds)

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
			handleObstruction(obstructionActivated, &elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer, controllerChannels, synchronizationChannels)

		case <-openDoorTimer.C:
			handleDoorTimeout(&elevator, openDoorTimer, orderChannels.ServicedOrderChannel, detectMotorFailureTimer)

		case <-detectMotorFailureTimer.C:
			handleDetectedMotorFailure(&elevator, detectMotorFailureTimer, controllerChannels, synchronizationChannels)

		case <-sendElevatorUpdateTicker.C:
			synchronizationChannels.LocalElevatorChannel <- elevator
		}
	}
}
