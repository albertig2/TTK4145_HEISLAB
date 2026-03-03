package main

import (
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	HallRequestsForAllIds := make(map[int][N_FLOORS][2]orderStatus)

	system := ElevatorSystem{
		HallRequests: [N_FLOORS][2]orderStatus{},
		States:       make(map[int]*ElevatorState),
	}
	initialize(&system, 1)
	initialize(&system, 2)
	initialize(&system, 3)
	setFloor(&system, 2, 2)
	setDirection(&system, 1, up)
	setBehavior(&system, 1, moving)
	//setCabRequests(&system, 1, 2, true)
	setHallRequests(&system, 2, hallUp, pending)

	HallRequestsForAllIds[1] = system.HallRequests
	HallRequestsForAllIds[2] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}
	HallRequestsForAllIds[3] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}

	hallRequestAssigner(&system, HallRequestsForAllIds, []int{1, 2, 3})
	t.Logf("system json: %s", encodeElevatorSystem(&system))
	transition := CheckOrderTransitionStatusForElevators(HallRequestsForAllIds, hallUp, 2, []int{1, 2, 3})
	t.Logf("transition: %v", transition)

	//message := encodeElevatorSystem(&system)
	//fmt.Print(message)
	//system2 := decodeElevatorSystem(message)
	//initialize(&system2, 4)
	//setCabRequests(&system2, 4, 1, true)

	//t.Logf("system2 json: %s", toJsonString(system2))

}
