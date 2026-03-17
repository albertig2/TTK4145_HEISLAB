package elevatorController

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

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

func initializeElevatorHardware(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer) {
	turnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)

	if FloorSensor() == -1 {
		initializeElevatorBetweenFloors(elevator, detectMotorFailureTimer)
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


func handleUpdatePeerViewOfCleardOrders(cleardOrderList []elevatorConfig.ButtonEvent, servicedOrderChannel chan elevatorConfig.ButtonEvent) {
	for _, cleardOrder := range cleardOrderList {
		servicedOrderChannel <- cleardOrder
	}
}

func handleOpenDoor(elevator *elevatorConfig.Elevator, openDoorTimer *time.Timer, timeOpenSeconds time.Duration, servicedOrderChannel chan elevatorConfig.ButtonEvent) {
	clearedOrders := []elevatorConfig.ButtonEvent{}

	if (FloorSensor() != -1) {
		fmt.Println("Door Open")

		elevio.SetDoorOpenLamp(true)
		elevator.Behavior = elevatorConfig.DoorOpen

		openDoorTimer.Stop()
		openDoorTimer.Reset(timeOpenSeconds)

		*elevator, clearedOrders = clearOrdersAtCurrentFloor(*elevator)
		clearListOfOrderLighst(clearedOrders)

		handleUpdatePeerViewOfCleardOrders(clearedOrders, servicedOrderChannel)
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

func handleOrderButtonPressd(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, buttonFloor int, buttonType elevatorConfig.Button, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
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

	tunrOnOrderLightsBasedOnLocalQueue(*elevator)

	fmt.Printf("\nNew state after HandleOrderButtonPres:\n")
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
			initializeElevatorHardware(elevator, detectMotorFailureTimer)

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
