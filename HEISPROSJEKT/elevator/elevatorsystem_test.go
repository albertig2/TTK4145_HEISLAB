package main

import (
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	system := ElevatorSystem{
		HallRequests: [N_FLOORS][2]bool{},
		States:       make(map[int]*ElevatorState),
	}

	initialize(&system, 1)
	initialize(&system, 2)
	initialize(&system, 3)
	setFloor(&system, 1, 2)
	setDirection(&system, 1, up)
	setBehavior(&system, 1, moving)
	//setCabRequests(&system, 1, 2, true)
	setHallRequests(&system, 2, true, true)
	hallRequestAssigner(&system)

	t.Logf("system json: %s", encodeElevatorSystem(&system))
	//message := encodeElevatorSystem(&system)
	//fmt.Print(message)
	//system2 := decodeElevatorSystem(message)
	//initialize(&system2, 4)
	//setCabRequests(&system2, 4, 1, true)

	//t.Logf("system2 json: %s", toJsonString(system2))

}
