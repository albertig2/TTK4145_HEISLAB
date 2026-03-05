package main

import (
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	HallRequestsForAllElevators := make(map[int][N_FLOORS][2]orderStatus)
	//CabRequestsForAllElevators := make(map[int][N_FLOORS]orderStatus)

	system := ElevatorSystem{}
	initialize(&system, 1)
	addPeer(&system, 2)
	addPeer(&system, 3)
	setFloor(&system, 2)
	setDirection(&system, up)
	setBehavior(&system, moving)
	//setCabRequests(&system, 1, 2, true)
	setHallRequests(&system, 2, hallUp, pending)

	HallRequestsForAllElevators[1] = system.HallRequests
	HallRequestsForAllElevators[2] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}
	HallRequestsForAllElevators[3] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}

	hallRequestAssigner(&system, HallRequestsForAllElevators, []int{1, 2, 3})
	t.Logf("system json: %s", encodeElevatorSystem(&system))
	transition := CheckOrderTransitionStatusForHallRequests(&system, HallRequestsForAllElevators, hallUp, 2, []int{1, 2, 3})
	t.Logf("transition: %v", transition)

	//message := encodeElevatorSystem(&system)
	//fmt.Print(message)
	//system2 := decodeElevatorSystem(message)
	//initialize(&system2, 4)
	//setCabRequests(&system2, 4, 1, true)

	//t.Logf("system2 json: %s", toJsonString(system2))

}
