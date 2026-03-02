package main

//"flag"

//flag.Int("id", 1, "Input id")

type Behavior int

const (
	idle     Behavior = 0
	moving   Behavior = 1
	doorOpen Behavior = 2
)

type ElevatorState struct {
	Behavior    Behavior
	Floor       int
	Direction   Dirn
	CabRequests [N_FLOORS]bool
}

type ElevatorSystem struct {
	HallRequests [N_FLOORS][2]bool
	States       map[int]*ElevatorState
}

func setBehavior(system *ElevatorSystem, id int, b Behavior) {
	state := system.States[id]
	state.Behavior = b
}

func setFloor(system *ElevatorSystem, id int, f int) {
	state := system.States[id]
	state.Floor = f
}

func setDirection(system *ElevatorSystem, id int, dir Dirn) {
	state := system.States[id]
	state.Direction = dir
}

// Usikker på om jeg skal kalle det on eller off? eller en funksjon for på og en for av
func setCabRequests(system *ElevatorSystem, id int, f int, on bool) {
	state := system.States[id]
	state.CabRequests[f] = on
}

// Usikker på om jeg skal kalle den up or down eller dont know
func setHallRequests(system *ElevatorSystem, f int, up bool, on bool) {
	if up {
		system.HallRequests[f][0] = on
	} else {
		system.HallRequests[f][1] = on
	}
}

func initialize(system *ElevatorSystem, id int) {
	// To decide floor can just do the get_floor_sensor_signal() and initialize to that floor, but for now hardcoded
	currentFloor := 1
	system.States[id] = &ElevatorState{
		Behavior:    idle,
		Floor:       currentFloor,
		Direction:   D_Stop,
		CabRequests: [N_FLOORS]bool{},
	}
}

// Spesify IDs as arguments when initializing (legge til et eller annet sted? Peer place??)

// Initialize, json, cost_function

// Funksjon for transisjoner mellom states (når man skal gå fra en state til en annen)
// Denne burde si hvilke transisjoner som skal gjøres (en pure function)
// Og så burde man ha en funskjon som utfører transjosjonen med påfølgende handlinger
// Union funksjon for å sette alle states til det andre sender som states
