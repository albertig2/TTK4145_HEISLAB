package elevatorHardware

import (
	"Driver-go/elevio"
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/timer"
	"fmt"
	"time"
	//"fmt"
	// "HEISPROSJEKT/elevatorConfig"
	// "HEISPROSJEKT/timer"
	// "fmt"
)

func SetAllLights(elevator elevatorConfig.Elevator) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			ElevatorRequestButtonLight(floor, elevatorConfig.Button(btn), elevator.Requests[floor][btn])
		}
	}
}
func TurnOffAllOrderLights() {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for button := elevio.ButtonType(0); button < 3; button++ {
			elevio.SetButtonLamp(button, floor, false)
		}
	}
}

func InitElevatorBetweenFloors(elevator *elevatorConfig.Elevator) {
	ElevatorMotorDirection(elevatorConfig.Down)
	elevator.Direction = elevatorConfig.Down
	elevator.Behavior = elevatorConfig.Moving
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
	}
	return hardwareChannels
}

// small initialisation sequence to put the elevator in a known behavior
// small initialisation sequence to put elevator in a known state
func InitElevatorHardware(elevator *elevatorConfig.Elevator) {

	TurnOffAllOrderLights()
	elevio.SetDoorOpenLamp(false)
	elevio.SetStopLamp(false)
	InitElevatorBetweenFloors(elevator)

}

func updateMotorDirection(motorDirection chan elevatorConfig.Direction) {
	for {
		d := <-motorDirection

		println("Motordirection:", d)

		elevio.SetMotorDirection(elevio.MotorDirection(d))
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

func HandleOnFloorArrival(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, newFloor int, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	Elevator_print(*elevator)

	elevator.Floor = newFloor
	ElevatorFloorIndicator(elevator.Floor)

	switch elevator.Behavior {
	case elevatorConfig.Moving:
		if Requests_shouldStop(*elevator) {

			ElevatorMotorDirection(elevatorConfig.Stop)

			//MotorDirectionChannel <- elevio.MD_Stop
			OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)
			//Elevator_doorLight(true)
			*elevator = Requests_clearAtCurrentFloor(*elevator)
			//timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
			//SetAllLights(*elevator)
			//elevator.Behavior = elevatorConfig.DoorOpen
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*elevator)
}

func OpenDoor(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, timeOpenSeconds time.Duration) {

	if !isBetweenFloors() {

		fmt.Println("Door Open")
		doorTimer.Stop()
		doorTimer.Reset(timeOpenSeconds)
		elevio.SetDoorOpenLamp(true)
		SetAllLights(*elevator)
		elevator.Behavior = elevatorConfig.DoorOpen

	} else {
		fmt.Println("Door was attempted opend in between floors")
		//Should probably trigger some sort of error handelig or somthing here
	}
}

func OnDoorTimeout(elevator *elevatorConfig.Elevator) {
	fmt.Printf("Door timeout")
	Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		pair := Requests_chooseDirection(*elevator)
		elevator.Direction = pair.Direction
		elevator.Behavior = pair.behavior

		switch elevator.Behavior {
		case elevatorConfig.DoorOpen:
			timer.Timer_start(elevator.Config.DoorOpenDuration_s)
			*elevator = Requests_clearAtCurrentFloor(*elevator)
			SetAllLights(*elevator)

		case elevatorConfig.Moving, elevatorConfig.Idle:
			ElevatorDoorLight(false)
			ElevatorMotorDirection(elevator.Direction)
		}
	default:
		// nothing
	}

	fmt.Printf("\nNew state:\n")
	Elevator_print(*elevator)
}

func HandleRequestButtonPress(elevator *elevatorConfig.Elevator, btn_floor int, btn_type elevatorConfig.Button, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Printf("\n\n%s(%d, %s)\n", "Recieved the following order:", btn_floor, elevatorConfig.ButtonToString(btn_type))
	Elevator_print(*elevator)

	switch elevator.Behavior {
	case elevatorConfig.DoorOpen:
		if Requests_shouldClearImmediately(*elevator, btn_floor, btn_type) {
			timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
		} else {
			elevator.Requests[btn_floor][btn_type] = true
		}

	case elevatorConfig.Moving:
		elevator.Requests[btn_floor][btn_type] = true

	case elevatorConfig.Idle:
		elevator.Requests[btn_floor][btn_type] = true
		pair := Requests_chooseDirection(*elevator)
		elevator.Direction = pair.Direction
		elevator.Behavior = pair.behavior

		switch pair.behavior {

		case elevatorConfig.DoorOpen:
			ElevatorDoorLight(true)
			timer.Timer_start(elevatorConfig.DOOR_OPEN_DURATION_S)
			*elevator = Requests_clearAtCurrentFloor(*elevator)

		case elevatorConfig.Moving:
			ElevatorMotorDirection(elevator.Direction)

		case elevatorConfig.Idle:
			// Do nothing
		}
	}

	SetAllLights(*elevator)

	fmt.Printf("\nNew state:\n")
	Elevator_print(*elevator)
}

func HandleStopButtonActivated (elevator *elevatorConfig.Elevator, doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Println("Stop was activated")
	TurnOffAllOrderLights()
	elevio.SetStopLamp(true)

	//MotorDirectionChannel <- elevatorConfig.Stop

	switch elevator.Behavior {

	case elevatorConfig.DoorOpen:
		//If stop is triggerd while at a floor, the door is opend and keep open until stop is reset + 3 seconds more
		OpenDoor(elevator, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)
		doorTimer.Stop()

	case elevatorConfig.Moving:
		ElevatorMotorDirection(elevatorConfig.Stop)

	default:

	}

}

func HandleStopButtonreset (elevator *elevatorConfig.Elevator, doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {
	fmt.Println("Stop was activated")
	TurnOffAllOrderLights()
	elevio.SetStopLamp(true)

	//MotorDirectionChannel <- elevatorConfig.Stop

	switch elevator.Behavior {

	case elevatorConfig.Moving,elevatorConfig.Idle:
		InitElevatorBetweenFloors(elevator)
		
	default:

	}

}

func HandleObstructionActivated(elevator *elevatorConfig.Elevator, doorTimer *time.Timer, MotorDirectionChannel chan elevatorConfig.Direction) {

	fmt.Println("Obstruction was activated")
	//ElevatorMotorDirection(elevatorConfig.Stop)
	//MotorDirectionChannel <- elevatorConfig.Stop

	switch elevator.Behavior {

	case elevatorConfig.DoorOpen:
		doorTimer.Stop()

	default:
		//Do nothing if the door is not open
	}

}

func RunElevatorHardware(elevatorID string, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt) {
	doorTimer := time.NewTimer(elevatorConfig.DOOR_OPEN_DURATION_S)
	doorTimer.Stop()
	elevatorObject := InitializeElevatorObject(elevatorID)

	InitElevatorHardware(&elevatorObject)

	//go updateMotorDirection(hardwareChannels.MotorDirectionChannel)

	for {
		select {
		case floor := <-hardwareChannels.FloorSensorChannel:

			HandleOnFloorArrival(&elevatorObject, doorTimer, floor, hardwareChannels.MotorDirectionChannel)

		case recievedOrder := <-hardwareChannels.PollOrderButtonsChannel:

			HandleRequestButtonPress(&elevatorObject, int(recievedOrder.Floor), elevatorConfig.Button(recievedOrder.Button), hardwareChannels.MotorDirectionChannel)

		case stopActivated := <-hardwareChannels.PollStopButtonChannel:
			if stopActivated {
				HandleStopButtonActivated(&elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel)
			} else {
				HandleStopButtonreset(&elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel)
			}

		case obstructionActivated := <-hardwareChannels.PollObstructionChannel:
			if obstructionActivated {
				HandleObstructionActivated(&elevatorObject, doorTimer, hardwareChannels.MotorDirectionChannel)
			}else {
				OpenDoor(&elevatorObject, doorTimer, elevatorConfig.DOOR_OPEN_DURATION_S)
			}

		case <-doorTimer.C:
			OnDoorTimeout(&elevatorObject)

		}

		select {
		case hardwareChannels.ElevatorObjectChannel <- elevatorObject:
		default:
		}
	}
}
