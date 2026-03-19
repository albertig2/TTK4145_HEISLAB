package elevatorController

import (
	"Driver-go/elevio"
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"time"
)

/*
The event handler contains functions to set the state machine output and decide the elevators’ reaction to the events that
are generated either locally or externally. In this project events are defined loosely and refer to inputs that warrant some
sort of calculation or action from the elevator. In addition to the  main handlers, the file also contains utility 
functions used by the handlers. 
*/

func InitializeControllerChannels() elevatorConfig.ControllerChannels {
	controllerChannels := elevatorConfig.ControllerChannels{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		PollFloorSensorChannel:  make(chan int),
		RestartElevatorChannel:  make(chan bool),
	}
	return controllerChannels
}

func initializeElevatorBetweenFloors(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer) {
	motorDirection(elevatorConfig.Down, detectMotorFailureTimer)
	elevator.Direction = elevatorConfig.Down
	elevator.Behavior = elevatorConfig.Moving
}

func initializeEmptyElevator(ownId string) elevatorConfig.Elevator {
	elevator := elevatorConfig.Elevator{
		OwnId:           ownId,
		Floor:           elevio.GetFloor(),
		Direction:       elevatorConfig.Stop,
		LocalOrderQueue: [elevatorConfig.NumberOfFloors][elevatorConfig.NunberOfButtons]bool{},
		Behavior:        elevatorConfig.Idle,
	}
	return elevator
}

func initializeElevatorHardware(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer) {
	turnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	stopButtonLight(false)

	if FloorSensor() == -1 {
		initializeElevatorBetweenFloors(elevator, detectMotorFailureTimer)
	}
}

func handleOnFloorArrival(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, newFloor int, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	elevator.Floor = newFloor
	floorIndicatorLight(elevator.Floor)

	switch elevator.Behavior {
	case elevatorConfig.Moving:
		if shouldStopAtCurrentFloor(*elevator) {
			motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)
		}
	default:
		// nothing
	}
}

func handleUpdatePeerViewOfCleardOrders(cleardOrderList []elevatorConfig.ButtonEvent, servicedOrderChannel chan elevatorConfig.ButtonEvent) {
	for _, cleardOrder := range cleardOrderList {
		servicedOrderChannel <- cleardOrder
	}
}

func handleOpenDoor(elevator *elevatorConfig.Elevator, openDoorTimer *time.Timer, timeOpenSeconds time.Duration, servicedOrderChannel chan elevatorConfig.ButtonEvent) {
	clearedOrders := []elevatorConfig.ButtonEvent{}

	if FloorSensor() != -1 {
		elevio.SetDoorOpenLamp(true)
		elevator.Behavior = elevatorConfig.DoorOpen

		openDoorTimer.Stop()
		openDoorTimer.Reset(timeOpenSeconds)

		*elevator, clearedOrders = clearOrdersAtCurrentFloor(*elevator)
		clearListOfOrderLighst(clearedOrders)

		handleUpdatePeerViewOfCleardOrders(clearedOrders, servicedOrderChannel)
	} else {
		//safety guard, do not open the door if in between floors
	}
}

func handleDoorTimeout(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		pair := chooseDirectionBasedOnOrders(*elevator)
		elevator.Direction = pair.direction
		elevator.Behavior = pair.behavior

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			doorLight(false)
			motorDirection(elevator.Direction, detectMotorFailureTimer)
		}
	default:
		// Do nothing
	}
}

func handleOrderButtonPressd(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, buttonFloor int, buttonType elevatorConfig.Button, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		if shouldClearOrderImmediately(*elevator, buttonFloor, buttonType) {
			elevator.LocalOrderQueue[buttonFloor][buttonType] = true
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)

		} else {
			elevator.LocalOrderQueue[buttonFloor][buttonType] = true
		}

	case elevatorConfig.Moving:
		elevator.LocalOrderQueue[buttonFloor][buttonType] = true

	case elevatorConfig.Idle:
		elevator.LocalOrderQueue[buttonFloor][buttonType] = true
		pair := chooseDirectionBasedOnOrders(*elevator)
		elevator.Direction = pair.direction
		elevator.Behavior = pair.behavior

		switch pair.behavior {
		case elevatorConfig.DoorOpen:
			elevator.LocalOrderQueue[buttonFloor][buttonType] = true
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)

		case elevatorConfig.Moving:
			motorDirection(elevator.Direction, detectMotorFailureTimer)

		case elevatorConfig.Idle:
			// Do nothing
		}
	}
	tunrOnOrderLightsBasedOnLocalQueue(*elevator)
}

func handleStopButton(stopActivated bool, elevator *elevatorConfig.Elevator, openDoorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {

	if stopActivated {
		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, openDoorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)
			openDoorTimer.Stop()

		case elevatorConfig.Moving:
			motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

		default:

		}
		turnOffAllOrderLights()
		stopButtonLight(true)

	} else {
		stopButtonLight(false)

		switch elevator.Behavior {
		case elevatorConfig.Moving, elevatorConfig.Idle:
			initializeElevatorHardware(elevator, detectMotorFailureTimer)

		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, openDoorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)
		default:
			//Do nothing
		}
	}
}

func handleObstruction(obstructionActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels) {

	if obstructionActivated {
		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			doorTimer.Stop()
			controllerChannels.RestartElevatorChannel <- true
			synchronisationChannels.PeerTransmitEnableChannel <- false
			turnOffAllOrderLights()
		default:
			//Do nothing
		}
	} else {
		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			handleRestartElevator(elevator, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DoorOpenDurationInSeconds, servicedOrderChannel)
		default:
			//Do nothing
		}

	}
}

func handleRestartElevator(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels) {

	motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

	synchronisationChannels.PeerTransmitEnableChannel <- false

	*elevator = initializeEmptyElevator(elevator.OwnId)

	controllerChannels.RestartElevatorChannel <- true

	initializeElevatorHardware(elevator, detectMotorFailureTimer)

	controllerChannels.RestartElevatorChannel <- false

	synchronisationChannels.PeerTransmitEnableChannel <- true
}

func handleLightSettingForPeerOrders(floor int, buttonType elevatorConfig.Button, lightValue bool) {
	elevio.SetButtonLamp(elevio.ButtonType(int(buttonType)), floor, lightValue)
}

func handleDetectedMotorFailure(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels) {

	synchronisationChannels.PeerTransmitEnableChannel <- false

	motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

	simulateMotorFailureTimer := time.NewTimer(2 * time.Second)

	<-simulateMotorFailureTimer.C

	handleRestartElevator(elevator, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
}
