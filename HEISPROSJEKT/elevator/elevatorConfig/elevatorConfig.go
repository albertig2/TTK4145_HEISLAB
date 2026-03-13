package elevatorConfig

import (
	"Driver-go/elevio"
	"time"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

const DOOR_OPEN_DURATION_S = 3*time.Second
const MOTOR_TIMEOUT_DURATION_S = 4*time.Second

type Direction int

const (
	Down Direction = -1
	Stop Direction = 0
	Up   Direction = 1
)

type Button int

const (
	HallUp   Button = 0
	HallDown Button = 1
	Cab      Button = 2
)

type ButtonEvent struct {
	Floor  int
	Button Button
}

type Behavior int

const (
	Idle     Behavior = 0
	DoorOpen Behavior = 1
	Moving   Behavior = 2
)

type Config struct {
	DoorOpenDuration_s time.Duration
	MotorTimeout_s     time.Duration
}

type Elevator struct {
	OwnId     string
	Floor     int
	Direction Direction
	Requests  [N_FLOORS][N_BUTTONS]bool
	Behavior  Behavior
	Config    Config
}

type OrderStatus string

// Cab order only needs to go from no order to Pending to Completed, while hall orders also need Assigned, since they are Assigned to an elevator by the assigner
const (
	Unknown   OrderStatus = "unknown"
	NoOrder   OrderStatus = "no order"
	Pending   OrderStatus = "pending"
	Assigned  OrderStatus = "assigned"
	Completed OrderStatus = "completed"
)

type ElevatorHardwareChannelsStruckt struct {
	PollOrderButtonsChannel chan elevio.ButtonEvent
	PollObstructionChannel  chan bool
	PollStopButtonChannel   chan bool
	FloorSensorChannel      chan int
	DoorOpenChannel         chan bool
	MotorDirectionChannel   chan Direction
	ElevatorObjectChannel   chan Elevator
	MotorFailureChannel     chan bool
}

type PeerRequestUpdate struct {
	PeerID   string
	HallReqs [N_FLOORS][2]OrderStatus
	CabReqs  [N_FLOORS]OrderStatus
}

type ElevatorOrderChannelStruckt struct {
	NewRecievedOrderChannel chan ButtonEvent //all new orders detected by the fsm from button presses are put here
	NewAssignedOrderChannel chan ButtonEvent //when a order is assigned to this elevator, the order is put here
	ServicedOrderChannel    chan ButtonEvent //when a order is serviced, it is put here
}

func DirectionToString(direction Direction) string {
	switch direction {
	case Up:
		return "up"
	case Down:
		return "down"
	case Stop:
		return "stop"
	default:
		return "undefined"
	}
}

func ButtonToString(b Button) string {
	switch b {
	case HallUp:
		return "hallUp"
	case HallDown:
		return "hallDown"
	case Cab:
		return "cab"
	default:
		return "undefined"
	}
}

func BehaviorToString(eb Behavior) string {
	switch eb {
	case Idle:
		return "idle"
	case DoorOpen:
		return "doorOpen"
	case Moving:
		return "moving"
	default:
		return "undefined"
	}
}

/*

type Config struct {
	doorOpenDuration_s time.Duration
}

type Elevator struct {
	floor     int
	direction	requests  [N_FLOORS][N_BUTTONS]bool
	behavior  ElevatorBehavior
	config    Config
}
*/
/*





 */
