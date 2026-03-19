package elevatorController

import (
	"Driver-go/elevio"
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"time"
)

/*
The event handler contains functions to set the state machine output and decide the elevators’
reaction to the events that are generated either locally or externally. In this project events are
defined loosely and refer to inputs that warrant some sort of calculation or action from the elevator.
The handlers are named after the event they correspond to. Each handler will set the correct output for
the elevator based on the input received, as well as change the state of the elevator in the finite state
machine. This includes all reactions to button presses, handling the assigned orders received from the order
assigner, ensuring correct light setting behavior and updating the order assigner of orders that has been
serviced.Most of the handlers are written on the finite state machine format, ensuring conditional response
based on the elevator behavior at the time the vent occurred. The exceptions are events that should be
handled the same regardless of previous state, like initializations or restarts.
There are ten different events who will trigger actions and outputs in the state machine.
Six of the events are local, stemming from hardware inputs or from the two timers. This includes
arriving at a floor, detecting a button press (order, stop or obstruction) or timeouts of the
door timer or, motor failure timer or the eleavtor update ticker. Three events are external, coming from other order assigner.
This includes a new order being assigned to this elevator, a new order assigned to a different peer
and a new order serviced by a different peer. Information about peer orders is necessary for the local
elevator to set and clear the correct order lights. The local panel should light up all orders being
handled by the system, not just the local orders.
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
		// nothing
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
			//do nothing
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
