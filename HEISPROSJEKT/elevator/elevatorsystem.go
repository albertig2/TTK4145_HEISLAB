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
	OwnId        int                      `json:"ownId"`
	HallRequests [N_FLOORS][2]orderStatus `json:"hallRequests"`
	States       map[int]*ElevatorState   `json:"states"`
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

func initialize(system *ElevatorSystem, id int) {
	// To decide floor can just do the get_floor_sensor_signal() and initialize to that floor, but for now hardcoded
	system.OwnId = id
	system.HallRequests = [N_FLOORS][2]orderStatus{}
	system.States = make(map[int]*ElevatorState)
	currentFloor := 1 // Get floor sensor. (men helst ikke -1? så siste faktisk floor)
	system.States[id] = &ElevatorState{
		Behavior:    idle,
		Floor:       currentFloor,
		Direction:   stop,
		CabRequests: [N_FLOORS]orderStatus{},
	}
}

func addPeer(system *ElevatorSystem, id int) {
	if _, exists := system.States[id]; !exists {
		system.States[id] = &ElevatorState{
			Behavior:    idle,
			Floor:       1, // Usikker på hva dette skal være (-1? for undefined until man får høre det fra heisen selv?)
			Direction:   stop,
			CabRequests: [N_FLOORS]orderStatus{},
		}
	}
}

// Spesify IDs as arguments when initializing (legge til et eller annet sted? Peer place??)

// Initialize, json, cost_function

// Funksjon for transisjoner mellom states (når man skal gå fra en state til en annen)
// Denne burde si hvilke transisjoner som skal gjøres (en pure function)
// Og så burde man ha en funskjon som utfører transjosjonen med påfølgende handlinger
// Union funksjon for å sette alle states til det andre sender som states

// Lage et map over alle ideene sine hallrequests, som brukes til å avgjøre state overganger.
// burde endre assigner til å bare endre sin egen hall request, men da med korrekt statet
// Man burde sikkert bare kunne sette floor osv på egen id og ikke på andres

// Kan hende man burde sende egen ID også under ElevatorState struct, lettere da å legge inn hallrequests riktig i henhold til den store matrisa alle holder på?
// Endre ider til å være strings i stedet for ints (matcher bedre med det Odin har gjort.
// Endre slik at man kun kan sette egne floors osv, og ikke andres, for å unngå feil
