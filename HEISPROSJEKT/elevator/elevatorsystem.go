package main

//"flag"

//flag.Int("id", 1, "Input id")

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

type ElevatorState struct {
	Behavior    Behavior       `json:"behaviour"`
	Floor       int            `json:"floor"`
	Direction   Direction      `json:"direction"`
	CabRequests [N_FLOORS]bool `json:"cabRequests"`
}

type ElevatorSystem struct {
	HallRequests [N_FLOORS][2]bool      `json:"hallRequests"`
	States       map[int]*ElevatorState `json:"states"`
}

func setBehavior(system *ElevatorSystem, id int, b Behavior) {
	state := system.States[id]
	state.Behavior = b
}

func setFloor(system *ElevatorSystem, id int, f int) {
	state := system.States[id]
	state.Floor = f
}

func setDirection(system *ElevatorSystem, id int, dir Direction) {
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
		Direction:   stop,
		CabRequests: [N_FLOORS]bool{},
	}
}

// Spesify IDs as arguments when initializing (legge til et eller annet sted? Peer place??)

// Initialize, json, cost_function

// Funksjon for transisjoner mellom states (når man skal gå fra en state til en annen)
// Denne burde si hvilke transisjoner som skal gjøres (en pure function)
// Og så burde man ha en funskjon som utfører transjosjonen med påfølgende handlinger
// Union funksjon for å sette alle states til det andre sender som states
