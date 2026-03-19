package elevatorController

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
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
		LocalOrderQueue: [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
		Behavior:        elevatorConfig.Idle,
	}
	return elevator
}

func InitializeControllerChannels() elevatorConfig.ElevatorControllerChannels {
	controllerChannels := elevatorConfig.ElevatorControllerChannels{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		FloorSensorChannel:      make(chan int),
		DoorOpenChannel:         make(chan bool),
		MotorDirectionChannel:   make(chan elevatorConfig.Direction),
		LocalElevatorChannel:    make(chan elevatorConfig.Elevator),
		RestartElevatorChannel:  make(chan bool),
	}
	return controllerChannels
}

func initializeElevatorHardware(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer) {
	turnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)

	if elevio.GetFloor() == -1 {
		initializeElevatorBetweenFloors(elevator, detectMotorFailureTimer)
	}
}

func isBetweenFloors() bool {
	currentFloor := elevio.GetFloor()
	if currentFloor != -1 {
		return false
	} else {
		return true
	}

}

func handleOnFloorArrival(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, newFloor int, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	fmt.Println("Elevator arrived at:", newFloor)
	debuggingHelpers.PrintLocalElvator(*elevator)

	elevator.Floor = newFloor
	floorIndicatorLight(elevator.Floor)

	switch elevator.Behavior {
	case elevatorConfig.Moving:
		if shouldStopAtCurrentFloor(*elevator) {
			motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)
		}
	default:
		// nothing
	}
	fmt.Printf("\nNew state after HandleOnFloorArrival:\n")
	debuggingHelpers.PrintLocalElvator(*elevator)
}

func handleOpenDoor(elevator *elevatorConfig.Elevator, openDoorTimer *time.Timer, timeOpenSeconds time.Duration, servicedOrderChannel chan elevatorConfig.ButtonEvent) {

	if !isBetweenFloors() {
		fmt.Println("Door Open")

		elevio.SetDoorOpenLamp(true)
		*elevator = clearOrdersAtCurrentFloor(*elevator, servicedOrderChannel)

		openDoorTimer.Stop()
		openDoorTimer.Reset(timeOpenSeconds)

		setAllOrderLights(*elevator)
		elevator.Behavior = elevatorConfig.DoorOpen
	} else {
		fmt.Println("Door was attempted opend in between floors")
	}
}

func handleDoorTimeout(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	fmt.Println("Door timeout")
	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		pair := chooseDirectionBasedOnOrders(*elevator)
		elevator.Direction = pair.direction
		elevator.Behavior = pair.behavior

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			doorLight(false)
			motorDirection(elevator.Direction, detectMotorFailureTimer)
		}
	default:
		// nothing
	}
	fmt.Printf("\nNew state after  OnDoorTimeout:\n")
	debuggingHelpers.PrintLocalElvator(*elevator)
}

func handleRequestButtonPressd(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, buttonFloor int, buttonType elevatorConfig.Button, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	fmt.Printf("\n\n%s(%d, %s)\n", "Recieved the following order:", buttonFloor, elevatorConfig.ButtonToString(buttonType))
	debuggingHelpers.PrintLocalElvator(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		if shouldClearOrderImmediately(*elevator, buttonFloor, buttonType) {
			elevator.LocalOrderQueue[buttonFloor][buttonType] = true
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)

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
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)

		case elevatorConfig.Moving:
			motorDirection(elevator.Direction, detectMotorFailureTimer)

		case elevatorConfig.Idle:
			// Do nothing
		}
	}

	setAllOrderLights(*elevator)

	fmt.Printf("\nNew state after HandleRequestButtonPres:\n")
	debuggingHelpers.PrintLocalElvator(*elevator)
}

func handleStopButton(stopActivated bool, elevator *elevatorConfig.Elevator, openDoorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {

	if stopActivated {
		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, openDoorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)
			openDoorTimer.Stop()

		case elevatorConfig.Moving:
			motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

		default:

		}
		fmt.Println("Stop was activated")
		turnOffAllOrderLights()
		elevio.SetStopLamp(true)

	} else {
		fmt.Println("Stop was reset")
		elevio.SetStopLamp(false)

		switch elevator.Behavior {
		case elevatorConfig.Moving, elevatorConfig.Idle:
			initializeElevatorBetweenFloors(elevator, detectMotorFailureTimer)

		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, openDoorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel) // keep door open for 3 more sek

		default:

		}

	}

}

func handleObstruction(obstructionActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ElevatorControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels) {

	if obstructionActivated {
		fmt.Println("Obstruction was activated")

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			doorTimer.Stop()
			controllerChannels.RestartElevatorChannel <- true
			synchronisationChannels.PeerTxEnableChannel <- false
		default:
			//Do nothing if the door is not open
		}

	} else {
		fmt.Println("Obstruction was reset")
		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			handleRestartElevator(elevator, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)
		default:
			//Do nothing if the door is not open
		}

	}
}

func handleRestartElevator(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ElevatorControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels) {

	motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

	synchronisationChannels.PeerTxEnableChannel <- false

	//make it unable to take new orders? (maybe not nessecerry)

	*elevator = initializeEmptyElevator(elevator.OwnId)

	controllerChannels.RestartElevatorChannel <- true

	initializeElevatorHardware(elevator, detectMotorFailureTimer)

	controllerChannels.RestartElevatorChannel <- false

	synchronisationChannels.PeerTxEnableChannel <- true

	fmt.Printf("\nNew state after motor failure detected:\n")
	debuggingHelpers.PrintLocalElvator(*elevator)
}

func handleLightSettingForPeerOrders(floor int, buttonType elevatorConfig.Button, lightValue bool) {
	elevio.SetButtonLamp(elevio.ButtonType(int(buttonType)), floor, lightValue)
}

func handleDetectedMotorFailure(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ElevatorControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels) {

	synchronisationChannels.PeerTxEnableChannel <- false

	motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

	simulateMotorFailureTimer := time.NewTimer(2 * time.Second)

	<-simulateMotorFailureTimer.C

	handleRestartElevator(elevator, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
}
