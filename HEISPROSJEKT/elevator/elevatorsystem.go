package main

//"flag"

//flag.Int("id", 1, "Input id")
type orderStatus string

// Cab order only needs to go from no order to pending to completed, while hall orders also need assigned, since they are assigned to an elevator by the assigner
const (
	noOrder   orderStatus = "no order"
	pending   orderStatus = "pending"
	assigned  orderStatus = "assigned"
	completed orderStatus = "completed"
)

type Behavior string

const (
	idle     Behavior = "idle"
	moving   Behavior = "moving"
	doorOpen Behavior = "doorOpen"
)

type Direction string

const (
	up   Direction = "up"
	down Direction = "down"
	stop Direction = "stop"
)

const (
	hallUp   = 0
	hallDown = 1
)

var HallDirs = [2]int{hallUp, hallDown}

type ElevatorState struct {
	Behavior    Behavior              `json:"behaviour"`
	Floor       int                   `json:"floor"`
	Direction   Direction             `json:"direction"`
	CabRequests [N_FLOORS]orderStatus `json:"cabRequests"`
}

type ElevatorSystem struct {
	OwnId        string                    `json:"ownId"`
	HallRequests [N_FLOORS][2]orderStatus  `json:"hallRequests"`
	States       map[string]*ElevatorState `json:"states"`
}

func setBehavior(system *ElevatorSystem, b Behavior) {
	state := system.States[system.OwnId]
	state.Behavior = b
}

func setFloor(system *ElevatorSystem, f int) {
	state := system.States[system.OwnId]
	state.Floor = f
}

func setDirection(system *ElevatorSystem, dir Direction) {
	state := system.States[system.OwnId]
	state.Direction = dir
}

// Usikker på om jeg skal kalle det on eller off? eller en funksjon for på og en for av
func setCabRequests(system *ElevatorSystem, f int, orderstatus orderStatus) {
	state := system.States[system.OwnId]
	state.CabRequests[f] = orderstatus
}

// Usikker på om jeg skal kalle den up or down eller dont know
func setHallRequests(system *ElevatorSystem, f int, halldir int, orderstatus orderStatus) {
	system.HallRequests[f][halldir] = orderstatus
}

func initialize(system *ElevatorSystem, id string) {
	// To decide floor can just do the get_floor_sensor_signal() and initialize to that floor, but for now hardcoded
	system.OwnId = id
	system.HallRequests = [N_FLOORS][2]orderStatus{}
	system.States = make(map[string]*ElevatorState)
	currentFloor := 1 // Get floor sensor. (men helst ikke -1? så siste faktisk floor)
	system.States[id] = &ElevatorState{
		Behavior:    idle,
		Floor:       currentFloor,
		Direction:   stop,
		CabRequests: [N_FLOORS]orderStatus{},
	}

	initializeHallRequests(system)
	initializeCabRequests(system)
}

func initializeHallRequests(system *ElevatorSystem) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for _, halldir := range HallDirs {
			system.HallRequests[floor][halldir] = noOrder
		}
	}
}

func initializeCabRequests(system *ElevatorSystem) {
	// Listen for other elevators to broadcast their view of your cab orders, and if you hear any, set your cab orders to be the combination of all of them (pending if any of them is pending or no order)
	// For each elevator you hear from, check all floors, if any of the floors have pending, set that floor to pending.
	for floor := 0; floor < N_FLOORS; floor++ {
		if system.States[system.OwnId].CabRequests[floor] != pending {
			system.States[system.OwnId].CabRequests[floor] = noOrder
		}
	}
}

// If only called from updatedElavatorSystemFromPeer, then I dont need the check for existence
func addPeer(system *ElevatorSystem, peerSystem *ElevatorSystem) {
	peerState := peerSystem.States[peerSystem.OwnId]
	if _, exists := system.States[peerSystem.OwnId]; !exists {
		system.States[peerSystem.OwnId] = &ElevatorState{
			Behavior:    peerState.Behavior,
			Floor:       peerState.Floor,
			Direction:   peerState.Direction,
			CabRequests: peerState.CabRequests,
		}
	}
}

func updateElevatorSystemFromPeer(system *ElevatorSystem, peerSystem *ElevatorSystem, HallRequestsForAllElevators map[string][N_FLOORS][2]orderStatus, CabRequestsForAllElevators map[string][N_FLOORS]orderStatus) {
	if _, exists := system.States[peerSystem.OwnId]; !exists {
		addPeer(system, peerSystem)
	}
	system.States[peerSystem.OwnId] = peerSystem.States[peerSystem.OwnId]

	HallRequestsForAllElevators[peerSystem.OwnId] = peerSystem.HallRequests
	if _, exists := peerSystem.States[system.OwnId]; exists {
		CabRequestsForAllElevators[peerSystem.OwnId] = peerSystem.States[system.OwnId].CabRequests
	}
}

// I utgangspunktet har jeg en annen funksjon som fikser andre transisjoner ...., kanskje nok å sette HallRequest listen to the appropriate, og så finne derfifra hva man skal sette
// Hvor ofte skal man sjekke transisjoner? med en gang etter man har updated

// Funksjon for transisjoner mellom states (når man skal gå fra en state til en annen)
// Denne burde si hvilke transisjoner som skal gjøres (en pure function)
// Og så burde man ha en funskjon som utfører transjosjonen med påfølgende handlinger

// Må deale med transisjoner, når man skal sette pending? når man skal gå til de andre? Skal man gjøre det når man får inn fra andre (hvertfall pending?)
// Når transisjon så man kansje gjøre ting også så jeg har jo en pure en
// Må på et tidspunkt oppdatere HallRequest med egen id sin hallrequests også og cabRequests.
