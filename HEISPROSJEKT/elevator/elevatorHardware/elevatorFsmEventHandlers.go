package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	"time"
)

func InitElevatorBetweenFloors(elevator *elevatorConfig.Elevator, motorTimeoutTimer *time.Timer) {
	ElevatorMotorDirection(elevatorConfig.Down, motorTimeoutTimer)
	elevator.Direction = elevatorConfig.Down
	elevator.Behavior = elevatorConfig.Moving
}

func InitializeElevator(id string) elevatorConfig.Elevator {

	config := elevatorConfig.Config{
		DoorOpenDuration_s: 3 * time.Second,
	}

	elevator := elevatorConfig.Elevator{
		OwnId:     id,
		Floor:     elevio.GetFloor(),
		Direction: elevatorConfig.Stop,
		Requests:  [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
		Behavior:  elevatorConfig.Idle,
		Config:    config,
	}

	return elevator
}

func InitElevatorHardwareChannels() elevatorConfig.ElevatorHardwareChannelsStruckt {

	hardwareChannels := elevatorConfig.ElevatorHardwareChannelsStruckt{
		PollOrderButtonsChannel: make(chan elevio.ButtonEvent),
		PollObstructionChannel:  make(chan bool),
		PollStopButtonChannel:   make(chan bool),
		FloorSensorChannel:      make(chan int),
		DoorOpenChannel:         make(chan bool),
		MotorDirectionChannel:   make(chan elevatorConfig.Direction),
		ElevatorObjectChannel:   make(chan elevatorConfig.Elevator),
		RestartElevatorChannel:  make(chan bool),
	}
	return hardwareChannels
}

// small initialisation sequence to put the elevator in a known behavior
// small initialisation sequence to put elevator in a known state
func InitElevatorHardware(elevator *elevatorConfig.Elevator, motorTimeoutTimer *time.Timer) {

	TurnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)

	if elevio.GetFloor() == -1 {
		InitElevatorBetweenFloors(elevator, motorTimeoutTimer)
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

func HandleOnFloorArrival(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, newFloor int, ServicedOrderChannel chan elevatorConfig.ButtonEvent, motorTimeoutTimer *time.Timer) {
	fmt.Println("Elevator arrived at:", newFloor)
	debuggingHelpers.Elevator_print(*elevator)

	elevator.Floor = newFloor
	ElevatorFloorIndicatorLight(elevator.Floor)

	switch elevator.Behavior {
	case elevatorConfig.Moving:
		if shouldStopAtCurrentFloor(*elevator) {

			ElevatorMotorDirection(elevatorConfig.Stop, motorTimeoutTimer)

			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state after HandleOnFloorArrival:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func OpenDoor(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, timeOpenSeconds time.Duration, ServicedOrderChannel chan elevatorConfig.ButtonEvent) {

	if !isBetweenFloors() {

		fmt.Println("Door Open")

		elevio.SetDoorOpenLamp(true)
		*elevator = clearOrdersAtCurrentFloor(*elevator, ServicedOrderChannel)

		doorTimer.Stop()
		doorTimer.Reset(timeOpenSeconds)

		SetAllLights(*elevator)
		elevator.Behavior = elevatorConfig.DoorOpen
		fmt.Println("Door Open end of if")

	} else {
		fmt.Println("Door was attempted opend in between floors")
		//Should probably trigger some sort of error handelig or somthing here
	}
}

func OnDoorTimeout(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, ServicedOrderChannel chan elevatorConfig.ButtonEvent, motorTimeoutTimer *time.Timer) {
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
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			ElevatorDoorLight(false)
			ElevatorMotorDirection(elevator.Direction, motorTimeoutTimer)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state after  OnDoorTimeout:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func HandleRequestButtonPress(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, btn_floor int, btn_type elevatorConfig.Button, ServicedOrderChannel chan elevatorConfig.ButtonEvent, motorTimeoutTimer *time.Timer) {
	fmt.Printf("\n\n%s(%d, %s)\n", "Recieved the following order:", btn_floor, elevatorConfig.ButtonToString(btn_type))
	debuggingHelpers.Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:

		if shouldClearOrderImmediately(*elevator, btn_floor, btn_type) {

			elevator.Requests[btn_floor][btn_type] = true
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel)

			//timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
		} else {
			elevator.Requests[btn_floor][btn_type] = true
		}

	case elevatorConfig.Moving:
		elevator.Requests[btn_floor][btn_type] = true

	case elevatorConfig.Idle:
		elevator.Requests[btn_floor][btn_type] = true
		pair := chooseDirectionBasedOnOrders(*elevator)
		elevator.Direction = pair.direction
		elevator.Behavior = pair.behavior

		switch pair.behavior {

		case elevatorConfig.DoorOpen:
			elevator.Requests[btn_floor][btn_type] = true
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel)

		case elevatorConfig.Moving:
			ElevatorMotorDirection(elevator.Direction, motorTimeoutTimer)

		case elevatorConfig.Idle:
			// Do nothing
		}
	}

	SetAllLights(*elevator)

	fmt.Printf("\nNew state after HandleRequestButtonPres:\n")
	debuggingHelpers.Elevator_print(*elevator)
}

func HandleStopButtonActivated(stopActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, ServicedOrderChannel chan elevatorConfig.ButtonEvent, motorTimeoutTimer *time.Timer) {

	if stopActivated {

		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel)
			doorTimer.Stop()

		case elevatorConfig.Moving:
			ElevatorMotorDirection(elevatorConfig.Stop, motorTimeoutTimer)

		default:

		}
		fmt.Println("Stop was activated")
		TurnOffAllOrderLights()
		elevio.SetStopLamp(true)

	} else {
		fmt.Println("Stop was reset")
		elevio.SetStopLamp(false)

		switch elevator.Behavior {

		case elevatorConfig.Moving, elevatorConfig.Idle:
			InitElevatorBetweenFloors(elevator, motorTimeoutTimer)

		case elevatorConfig.DoorOpen:
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel) // keep door open for 3 more sek

		default:

		}

	}

}

func HandleObstructionActivated(obstructionActivated bool, elevator *elevatorConfig.Elevator, doorTimer *time.Timer, ServicedOrderChannel chan elevatorConfig.ButtonEvent, motorTimeoutTimer *time.Timer, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels) {

	if obstructionActivated {
		fmt.Println("Obstruction was activated")

		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			doorTimer.Stop()
			hardwareChannels.RestartElevatorChannel <- true
			synchronisationChannels.PeerTxEnableChannel <- false

		default:
			//Do nothing if the door is not open
		}

	} else {
		fmt.Println("Obstruction was reset")

		switch elevator.Behavior {

		case elevatorConfig.DoorOpen:
			HandleRestartElevator(elevator, motorTimeoutTimer, hardwareChannels, synchronisationChannels)
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S, ServicedOrderChannel)
		default:
			//Do nothing if the door is not open
		}

	}

}

func HandleRestartElevator(elevatorObject *elevatorConfig.Elevator, motorTimeoutTimer *time.Timer, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels) {

	ElevatorMotorDirection(elevatorConfig.Stop, motorTimeoutTimer)

	synchronisationChannels.PeerTxEnableChannel <- false

	//make it unable to take new orders? (maybe not nessecerry)

	*elevatorObject = InitializeElevator(elevatorObject.OwnId)

	hardwareChannels.RestartElevatorChannel <- true

	InitElevatorHardware(elevatorObject, motorTimeoutTimer)

	hardwareChannels.RestartElevatorChannel <- false

	synchronisationChannels.PeerTxEnableChannel <- true

	fmt.Printf("\nNew state after motor failure detected:\n")
	debuggingHelpers.Elevator_print(*elevatorObject)
}
func HandlelightSettingForPeerOrders(floor int, buttonType elevatorConfig.Button, lightValue bool) {

	elevio.SetButtonLamp(elevio.ButtonType(int(buttonType)), floor, lightValue)

}

func HandleMotorTimeout(elevatorObject *elevatorConfig.Elevator, motorTimeoutTimer *time.Timer, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt, synchronisationChannels elevatorConfig.SynchronisationChannels) {
	synchronisationChannels.PeerTxEnableChannel <- false
	ElevatorMotorDirection(elevatorConfig.Stop, motorTimeoutTimer)
	simulateMotorFailureTimer := time.NewTimer(2 * time.Second)
	<-simulateMotorFailureTimer.C
	HandleRestartElevator(elevatorObject, motorTimeoutTimer, hardwareChannels, synchronisationChannels)
}
