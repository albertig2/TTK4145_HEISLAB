package elevatorConfig

import (
	"Driver-go/elevio"
	"Network-go/network/peers"
	"time"
)

const NumberOfFloors int = 4
const NunberOfButtons int = 3
const NumberOfHallButtons int = 2
const PeerUpdatePort int = 30004
const BroadcastPort int = 30400

const DoorOpenDurationInSeconds = 3 * time.Second
const MotorTimeOutDurationInSeconds = 4 * time.Second

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

type Elevator struct {
	OwnId           string
	Floor           int
	Direction       Direction
	LocalOrderQueue [NumberOfFloors][NunberOfButtons]bool
	Behavior        Behavior
}

type PeerState struct {
	Behavior  Behavior                    `json:"behavior"`
	Floor     int                         `json:"floor"`
	Direction Direction                   `json:"direction"`
	CabOrders [NumberOfFloors]OrderStatus `json:"cabOrders"`
}

type PeerView struct {
	AlivePeers []string                                         `json:"alivePeers"`
	OwnId      string                                           `json:"ownId"`
	HallOrders [NumberOfFloors][NumberOfHallButtons]OrderStatus `json:"hallOrders"`
	States     map[string]*PeerState                            `json:"states"`
}

type OrderStatus string

const (
	Unknown  OrderStatus = "unknown"
	NoOrder  OrderStatus = "no order"
	Pending  OrderStatus = "pending"
	Assigned OrderStatus = "assigned"
	Serviced OrderStatus = "serviced"
)

type SynchronizationChannels struct {
	PeerUpdateChannel                                  chan peers.PeerUpdate
	PeerTransmitEnableChannel                          chan bool
	BroadcastIncomingMessagesChannel                   chan PeerView
	BroadcastOutgoingMessagesChannel                   chan PeerView
	UpdatePeerViewWithLocalElevatorChannel             chan Elevator
	LocalElevatorChannel                               chan Elevator
	UpdatePeerViewforBroadcastWithLocalPeerViewChannel chan PeerView
	UpdateLocalPeerViewWithExternalPeerViewChannel     chan PeerView
	AlivePeersChannel                                  chan []string
}

type ControllerChannels struct {
	PollOrderButtonsChannel chan elevio.ButtonEvent
	PollObstructionChannel  chan bool
	PollStopButtonChannel   chan bool
	PollFloorSensorChannel  chan int
	RestartElevatorChannel  chan bool
}

type OrderChannels struct {
	NewRecievedOrderChannel     chan ButtonEvent
	NewAssignedOrderChannel     chan ButtonEvent
	NewAssignedPeerOrderChannel chan ButtonEvent
	ServicedOrderChannel        chan ButtonEvent
	ServicedPeerOrderChannel    chan ButtonEvent
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
