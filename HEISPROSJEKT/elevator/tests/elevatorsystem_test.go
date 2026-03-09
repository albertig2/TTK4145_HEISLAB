package elevatorHardware_test

import (
	"HEISPROSJEKT/communication"
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/orderProtocol"
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	// Sholuld always update these maps with own state before checking for transitions or doing anything based on the status of the floors
	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorHardware.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorHardware.OrderStatus)

	system1 := elevatorHardware.ElevatorSystem{}
	elevatorHardware.InitializeElevatorSystem(&system1, "1")
	elevatorHardware.SetFloor(&system1, 3)
	elevatorHardware.SetDirection(&system1, elevatorConfig.Stop)
	elevatorHardware.SetBehavior(&system1, elevatorConfig.DoorOpen)
	elevatorHardware.SetHallRequests(&system1, 2, int(elevatorConfig.HallDown), elevatorHardware.Pending)

	system2 := elevatorHardware.ElevatorSystem{}
	elevatorHardware.InitializeElevatorSystem(&system2, "2")
	elevatorHardware.SetFloor(&system2, 0)
	elevatorHardware.SetDirection(&system2, elevatorConfig.Stop)
	elevatorHardware.SetBehavior(&system2, elevatorConfig.Idle)
	elevatorHardware.SetHallRequests(&system2, 3, int(elevatorConfig.HallUp), elevatorHardware.Pending)

	system3 := elevatorHardware.ElevatorSystem{}
	elevatorHardware.InitializeElevatorSystem(&system3, "3")
	elevatorHardware.SetFloor(&system3, 1)
	elevatorHardware.SetDirection(&system3, elevatorConfig.Up)
	elevatorHardware.SetBehavior(&system3, elevatorConfig.Moving)
	elevatorHardware.SetCabRequests(&system3, 1, elevatorHardware.Pending)
	elevatorHardware.SetHallRequests(&system3, 3, int(elevatorConfig.HallUp), elevatorHardware.Pending)

	elevatorHardware.UpdateElevatorSystemFromPeer(&system1, &system2, HallRequestsForAllElevators, CabRequestsForAllElevators)
	elevatorHardware.UpdateElevatorSystemFromPeer(&system1, &system3, HallRequestsForAllElevators, CabRequestsForAllElevators)

	HallRequestsForAllElevators["1"] = system1.HallRequests           // Always update these before checking for transitions or doing anything based on the status of the floors
	CabRequestsForAllElevators["1"] = system1.States["1"].CabRequests // Always update these before checking for transitions or doing anything based on the status of the floors

	t.Logf("HallRequest for elevator 1: %v", HallRequestsForAllElevators["1"])
	orderTransition := orderProtocol.CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, int(elevatorConfig.HallUp), 3, []string{"1", "2", "3"})
	t.Logf("orderTransition: %v", orderTransition)
	elevatorHardware.SetHallRequests(&system1, 3, int(elevatorConfig.HallUp), elevatorHardware.Pending)
	HallRequestsForAllElevators["1"] = system1.HallRequests // Remember to update
	orderTransition2 := orderProtocol.CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, int(elevatorConfig.HallUp), 3, []string{"1", "2", "3"})
	t.Logf("orderTransition2: %v", orderTransition2)

	hallRequestTransitions := orderProtocol.GetAllHallRequestTransitions(&system1, HallRequestsForAllElevators, []string{"1", "2", "3"})
	orderProtocol.HallRequestAssigner(&system1, hallRequestTransitions, []string{"1", "2", "3"})
	t.Logf("HallRequest for system1 before assigning: %v", system1.HallRequests[3][int(elevatorConfig.HallUp)])
	orderProtocol.TransitionForHallRequestsByType(&system1, hallRequestTransitions, orderProtocol.PendingToAssigned, []string{"1", "2", "3"})
	t.Logf("HallRequest for system1 after assigning: %v", system1.HallRequests[3][int(elevatorConfig.HallUp)])

	t.Logf("system json: %s", communication.EncodeElevatorSystem(&system1))
	transition := orderProtocol.CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, int(elevatorConfig.HallUp), 2, []string{"1", "2", "3"})
	t.Logf("transition: %v", transition)
	t.Logf("HallRequests: %v", HallRequestsForAllElevators)
	t.Logf("CabRequests: %v", CabRequestsForAllElevators)
	//message := communication.EncodeElevatorSystem(&system)
	//fmt.Print(message)
	//system2 := communication.DecodeElevatorSystem(message)
	//orderProtocol.Initialize(&system2, 4)
	//setCabRequests(&system2, 4, 1, true)

	//t.Logf("system2 json: %s", toJsonString(system2))

}
