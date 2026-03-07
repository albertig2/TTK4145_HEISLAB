package elevatorHardware_test

import (
	"testing"
	"HEISPROSJEKT/orderProtocol"
	"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/communication"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	HallRequestsForAllElevators := make(map[string][elevatorHardware.N_FLOORS][2]elevatorHardware.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorHardware.N_FLOORS]elevatorHardware.OrderStatus)

	system1 := elevatorHardware.ElevatorSystem{}
	elevatorHardware.Initialize(&system1, "1")
	elevatorHardware.SetFloor(&system1, 2)
	elevatorHardware.SetDirection(&system1, elevatorHardware.Up)
	elevatorHardware.SetBehavior(&system1, elevatorHardware.Moving)
	elevatorHardware.SetCabRequests(&system1, 1, elevatorHardware.Pending)

	system2 := elevatorHardware.ElevatorSystem{}
	elevatorHardware.Initialize(&system2, "2")
	elevatorHardware.SetFloor(&system2, 3)
	elevatorHardware.SetDirection(&system2, elevatorHardware.Stop)
	elevatorHardware.SetBehavior(&system2, elevatorHardware.Idle)
	elevatorHardware.SetCabRequests(&system2, 2, elevatorHardware.Assigned)
	elevatorHardware.SetHallRequests(&system2, 3, elevatorHardware.HallUp, elevatorHardware.Pending)

	system3 := elevatorHardware.ElevatorSystem{}
	elevatorHardware.Initialize(&system3, "3")
	elevatorHardware.SetFloor(&system3, 3)
	elevatorHardware.SetDirection(&system3, elevatorHardware.Stop)
	elevatorHardware.SetBehavior(&system3, elevatorHardware.DoorOpen)
	elevatorHardware.SetHallRequests(&system3, 2, elevatorHardware.HallDown, elevatorHardware.Pending)
	elevatorHardware.SetHallRequests(&system3, 3, elevatorHardware.HallUp, elevatorHardware.Pending)

	elevatorHardware.UpdateElevatorSystemFromPeer(&system1, &system2, HallRequestsForAllElevators, CabRequestsForAllElevators)
	elevatorHardware.UpdateElevatorSystemFromPeer(&system1, &system3, HallRequestsForAllElevators, CabRequestsForAllElevators)

	HallRequestsForAllElevators["1"] = system1.HallRequests
	CabRequestsForAllElevators["1"] = system1.States["1"].CabRequests
	//HallRequestsForAllElevators["2"] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}
	//HallRequestsForAllElevators["3"] = [N_FLOORS][2]orderStatus{{noOrder, noOrder}, {noOrder, noOrder}, {pending, noOrder}, {noOrder, noOrder}}

	orderProtocol.HallRequestAssigner(&system1, HallRequestsForAllElevators, []string{"1", "2", "3"})
	t.Logf("system json: %s", communication.EncodeElevatorSystem(&system1))
	transition := orderProtocol.CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, elevatorHardware.HallUp, 2, []string{"1", "2", "3"})
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
