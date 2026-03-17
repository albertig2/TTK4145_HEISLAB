package elevatorController

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func initializeElevatorBetweenFloors(elevator *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer) {
	motorDirection(elevatorConfig.Down, detectMotorFailureTimer)
	elevator.Direction = elevatorConfig.Down
	elevator.Behavior = elevatorConfig.Moving
}

func initializeElevator(ownId string) elevatorConfig.Elevator {

	config := elevatorConfig.Config{
		DoorOpenDuration_s: 3 * time.Second,
	}

	elevator := elevatorConfig.Elevator{
		OwnId:     ownId,
		Floor:     elevio.GetFloor(),
		Direction: elevatorConfig.Stop,
		Requests:  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
		Behavior:  elevatorConfig.Idle,
		Config:    config,
	}

	return elevator
}

func InitializeControllerChannels() elevatorConfig.ElevatorHardwareChannelsStruckt {

	controllerChannels := elevatorConfig.ElevatorHardwareChannelsStruckt{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		FloorSensorChannel:      make(chan int),
		DoorOpenChannel:         make(chan bool),
		MotorDirectionChannel:   make(chan elevatorConfig.Direction),
		ElevatorObjectChannel:   make(chan elevatorConfig.Elevator),
		RestartElevatorChannel:  make(chan bool),
	}
	return controllerChannels
}

// small initialisation sequence to put the elevator in a known behavior
// small initialisation sequence to put elevator in a known state
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
	debuggingHelpers.Elevator_print(*elevator)

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
	debuggingHelpers.Elevator_print(*elevator)
}

func handleOpenDoor(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, timeOpenSeconds time.Duration, servicedOrderChannel chan elevatorConfig.ButtonEvent) {

	if !isBetweenFloors() {

		fmt.Println("Door Open")

		elevio.SetDoorOpenLamp(true)
		*elevator = clearOrdersAtCurrentFloor(*elevator, servicedOrderChannel)

		doorTimer.Stop()
		doorTimer.Reset(timeOpenSeconds)

		setAllOrderLights(*elevator)
		elevator.Behavior = elevatorConfig.DoorOpen
		fmt.Println("Door Open end of if")

	} else {
		fmt.Println("Door was attempted opend in between floors")
		//Should probably trigger some sort of error handelig or somthing here
	}
}

func handleDoorTimeout(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	fmt.Println("Door timeout")
	// debuggingHelpers.Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		pair := chooseDirectionBasedOnOrders(*elevator)
		elevator.Direction = pair.direction
		elevator.Behavior = pair.behavior

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			// timer.Timer_start(elevator.Config.DoorOpenDuration_s)
			// *elevator = RequestsClearAtCurrentFloor(*elevator)
			// SetAllLights(*elevator)
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			doorLight(false)
			motorDirection(elevator.Direction, detectMotorFailureTimer)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state after  OnDoorTimeout:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func handleRequestButtonPressd(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, buttonFloor int, buttonType elevatorConfig.Button, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {
	fmt.Printf("\n\n%s(%d, %s)\n", "Recieved the following order:", buttonFloor, elevatorConfig.ButtonToString(buttonType))
	debuggingHelpers.Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:

		if shouldClearOrderImmediately(*elevator, buttonFloor, buttonType) {

			elevator.Requests[buttonFloor][buttonType] = true
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)

			//timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
		} else {
			elevator.Requests[buttonFloor][buttonType] = true
		}

	case elevatorConfig.Moving:
		elevator.Requests[buttonFloor][buttonType] = true

	case elevatorConfig.Idle:
		elevator.Requests[buttonFloor][buttonType] = true
		pair := chooseDirectionBasedOnOrders(*elevator)
		elevator.Direction = pair.direction
		elevator.Behavior = pair.behavior

		switch pair.behavior {

		case elevatorConfig.DoorOpen:
			elevator.Requests[buttonFloor][buttonType] = true
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)

		case elevatorConfig.Moving:
			motorDirection(elevator.Direction, detectMotorFailureTimer)

		case elevatorConfig.Idle:
			// Do nothing
		}
	}

	setAllOrderLights(*elevator)

	fmt.Printf("\nNew state after HandleRequestButtonPres:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func handleStopButton(stopActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer) {

	if stopActivated {

		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel)
			doorTimer.Stop()

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
			handleOpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, servicedOrderChannel) // keep door open for 3 more sek

		default:

		}

	}

}

func handleObstruction(obstructionActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, servicedOrderChannel chan elevatorConfig.ButtonEvent, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels) {

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

func handleRestartElevator(elevatorObject *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels) {

	motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)

	synchronisationChannels.PeerTxEnableChannel <- false

	//make it unable to take new orders? (maybe not nessecerry)

	*elevatorObject = initializeElevator(elevatorObject.OwnId)

	controllerChannels.RestartElevatorChannel <- true

	initializeElevatorHardware(elevatorObject, detectMotorFailureTimer)

	controllerChannels.RestartElevatorChannel <- false

	synchronisationChannels.PeerTxEnableChannel <- true

	fmt.Printf("\nNew state after motor failure detected:\n")
	debuggingHelpers.Elevator_print(*elevatorObject)
}
func handleLightSettingForPeerOrders(floor int, buttonType elevatorConfig.Button, lightValue bool) {

	elevio.SetButtonLamp(elevio.ButtonType(int(buttonType)), floor, lightValue)

}

func handleDetectedMotorFailure(elevatorObject *elevatorConfig.Elevator, detectMotorFailureTimer *time.Timer, controllerChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels) {
	synchronisationChannels.PeerTxEnableChannel <- false
	motorDirection(elevatorConfig.Stop, detectMotorFailureTimer)
	simulateMotorFailureTimer := time.NewTimer(2 * time.Second)
	<-simulateMotorFailureTimer.C
	handleRestartElevator(elevatorObject, detectMotorFailureTimer, controllerChannels, synchronisationChannels)
}
