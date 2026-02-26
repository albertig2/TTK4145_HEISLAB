package main

import "testing"

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
	setDirection(&system, 1, D_Up)
	setBehavior(&system, 1, moving)
	setCabRequests(&system, 1, 2, true)
	setHallRequests(&system, 2, true, true)

	t.Logf("system json: %s", toJsonString(system))
}
