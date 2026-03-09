package elevatorsystem_test

import (
	"HEISPROSJEKT/synchronisation"
	"HEISPROSJEKT/elevatorConfig"
	//"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/orderProtocol"
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	// Sholuld always update these maps with own state before checking for transitions or doing anything based on the status of the floors
	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus)

	system1 := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system1, "1")
	synchronisation.SetFloor(&system1, 3)
	synchronisation.SetDirection(&system1, elevatorConfig.Stop)
	synchronisation.SetBehavior(&system1, elevatorConfig.DoorOpen)
	synchronisation.SetHallRequests(&system1, 2, int(elevatorConfig.HallDown), synchronisation.Pending)

	system2 := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system2, "2")
	synchronisation.SetFloor(&system2, 0)
	synchronisation.SetDirection(&system2, elevatorConfig.Stop)
	synchronisation.SetBehavior(&system2, elevatorConfig.Idle)
	synchronisation.SetHallRequests(&system2, 3, int(elevatorConfig.HallUp), synchronisation.Pending)

	system3 := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system3, "3")
	synchronisation.SetFloor(&system3, 1)
	synchronisation.SetDirection(&system3, elevatorConfig.Up)
	synchronisation.SetBehavior(&system3, elevatorConfig.Moving)
	synchronisation.SetCabRequests(&system3, 1, synchronisation.Pending)
	synchronisation.SetHallRequests(&system3, 3, int(elevatorConfig.HallUp), synchronisation.Pending)

	synchronisation.UpdateElevatorSystemFromPeer(&system1, &system2, HallRequestsForAllElevators, CabRequestsForAllElevators)
	synchronisation.UpdateElevatorSystemFromPeer(&system1, &system3, HallRequestsForAllElevators, CabRequestsForAllElevators)

	HallRequestsForAllElevators["1"] = system1.HallRequests           // Always update these before checking for transitions or doing anything based on the status of the floors
	CabRequestsForAllElevators["1"] = system1.States["1"].CabRequests // Always update these before checking for transitions or doing anything based on the status of the floors

	t.Logf("HallRequest for elevator 1: %v", HallRequestsForAllElevators["1"])
	orderTransition := orderProtocol.CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, int(elevatorConfig.HallUp), 3, []string{"1", "2", "3"})
	t.Logf("orderTransition: %v", orderTransition)
	synchronisation.SetHallRequests(&system1, 3, int(elevatorConfig.HallUp), synchronisation.Pending)
	HallRequestsForAllElevators["1"] = system1.HallRequests // Remember to update
	orderTransition2 := orderProtocol.CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, int(elevatorConfig.HallUp), 3, []string{"1", "2", "3"})
	t.Logf("orderTransition2: %v", orderTransition2)

	hallRequestTransitions := orderProtocol.GetAllHallRequestTransitions(&system1, HallRequestsForAllElevators, []string{"1", "2", "3"})
	orderProtocol.HallRequestAssigner(&system1, hallRequestTransitions, []string{"1", "2", "3"})
	t.Logf("HallRequest for system1 before assigning: %v", system1.HallRequests[3][int(elevatorConfig.HallUp)])
	orderProtocol.TransitionForHallRequestsByType(&system1, hallRequestTransitions, orderProtocol.PendingToAssigned, []string{"1", "2", "3"})
	t.Logf("HallRequest for system1 after assigning: %v", system1.HallRequests[3][int(elevatorConfig.HallUp)])

	t.Logf("system json: %s", synchronisation.EncodeElevatorSystem(&system1))
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
