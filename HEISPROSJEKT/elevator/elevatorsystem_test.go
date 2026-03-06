package main

import (
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	HallRequestsForAllElevators := make(map[string][N_FLOORS][2]orderStatus)
	CabRequestsForAllElevators := make(map[string][N_FLOORS]orderStatus)

	system1 := ElevatorSystem{}
	initialize(&system1, "1")
	setFloor(&system1, 2)
	setDirection(&system1, up)
	setBehavior(&system1, moving)
	setCabRequests(&system1, 1, pending)

	system2 := ElevatorSystem{}
	initialize(&system2, "2")
	setFloor(&system2, 3)
	setDirection(&system2, stop)
	setBehavior(&system2, idle)
	setCabRequests(&system2, 2, assigned)
	setHallRequests(&system2, 3, hallUp, pending)

	system3 := ElevatorSystem{}
	initialize(&system3, "3")
	setFloor(&system3, 3)
	setDirection(&system3, stop)
	setBehavior(&system3, doorOpen)
	setHallRequests(&system3, 2, hallDown, pending)
	setHallRequests(&system3, 3, hallUp, pending)

	updateElevatorSystemFromPeer(&system1, &system2, HallRequestsForAllElevators, CabRequestsForAllElevators)
	updateElevatorSystemFromPeer(&system1, &system3, HallRequestsForAllElevators, CabRequestsForAllElevators)

	//HallRequestsForAllElevators["1"] = system1.HallRequests
	//HallRequestsForAllElevators["2"] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}
	//HallRequestsForAllElevators["3"] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}

	hallRequestAssigner(&system1, HallRequestsForAllElevators, []string{"1", "2", "3"})
	t.Logf("system json: %s", encodeElevatorSystem(&system1))
	transition := CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, hallUp, 2, []string{"1", "2", "3"})
	t.Logf("transition: %v", transition)
	t.Logf("HallRequests: %v", HallRequestsForAllElevators)
	t.Logf("CabRequests: %v", CabRequestsForAllElevators)
	//message := encodeElevatorSystem(&system)
	//fmt.Print(message)
	//system2 := decodeElevatorSystem(message)
	//initialize(&system2, 4)
	//setCabRequests(&system2, 4, 1, true)

	//t.Logf("system2 json: %s", toJsonString(system2))

}
