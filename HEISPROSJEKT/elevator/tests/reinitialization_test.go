package elevatorsystem_test

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/orderProtocol"
	"HEISPROSJEKT/synchronisation"
	"testing"
)

// go test -v -run TestReinitializeCabRequestsRecovery

func TestReinitializeCabRequestsRecovery(t *testing.T) {
	//
	HallRequestsForAllElevators1 := make(map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus)
	CabRequestsForAllElevators1 := make(map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus)

	// Step 1: Setup initial system with a cab request
	system := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system, "1")

	system3 := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system3, "3")
	HallRequestsForAllElevators3 := make(map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus)
	CabRequestsForAllElevators3 := make(map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus)

	// Simulate peer's view before restart
	peerSystem := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&peerSystem, "2")
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		peerSystem.States[peerSystem.OwnId].CabRequests[floor] = synchronisation.NoOrder
	}
	synchronisation.SetCabRequests(&peerSystem, 2, synchronisation.Pending)

	synchronisation.UpdateElevatorSystemWithPeer(&system3, &peerSystem, HallRequestsForAllElevators3, CabRequestsForAllElevators3)

	synchronisation.SetCabRequests(&peerSystem, 3, synchronisation.Pending)

	synchronisation.UpdateElevatorSystemWithPeer(&system, &peerSystem, HallRequestsForAllElevators1, CabRequestsForAllElevators1)
	t.Logf("System cab requests before reinitialization: %v", system.States[peerSystem.OwnId].CabRequests)
	HallRequestsForAllElevators1[system.OwnId] = system.HallRequests
	CabRequestsForAllElevators1[system.OwnId] = system.States[system.OwnId].CabRequests

	t.Logf("Peer system cab requests before reinitialization; %v", peerSystem.States[peerSystem.OwnId].CabRequests)
	// Step 2: Simulate restart (cab requests set to Unknown)
	synchronisation.InitializeElevatorSystem(&peerSystem, "2")

	t.Logf("Peer system cab requests after reinitialization: %v", peerSystem.States[peerSystem.OwnId].CabRequests)
	HallRequestsForAllElevators2 := make(map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus)
	CabRequestsForAllElevators2 := make(map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus)

	synchronisation.UpdateElevatorSystemWithPeer(&peerSystem, &system, HallRequestsForAllElevators2, CabRequestsForAllElevators2)
	synchronisation.UpdateElevatorSystemWithPeer(&peerSystem, &system3, HallRequestsForAllElevators2, CabRequestsForAllElevators2)

	HallRequestsForAllElevators2[peerSystem.OwnId] = peerSystem.HallRequests
	CabRequestsForAllElevators2[peerSystem.OwnId] = peerSystem.States[peerSystem.OwnId].CabRequests
	t.Logf("Peer system cab requests after updating with peer: %v", peerSystem.States[peerSystem.OwnId].CabRequests)

	synchronisation.UpdateElevatorSystemWithPeer(&system, &peerSystem, HallRequestsForAllElevators1, CabRequestsForAllElevators1)
	t.Logf("System cab requests after updating with peer: %v", system.States[peerSystem.OwnId].CabRequests)

	// Step 3: Simulate receiving peer state

	alivePeers := []string{"1", "2", "3"}

	t.Logf("CabRequests matrix: %v", CabRequestsForAllElevators2)

	// Step 4: Run transition logic for cab requests
	transitions := orderProtocol.GetAllCabRequestTransitions(&peerSystem, CabRequestsForAllElevators2, alivePeers)
	t.Logf("Cab request transitions: %v", transitions)
	orderProtocol.TransitionForCabRequestsByType(&peerSystem, transitions, orderProtocol.UnknownToPending)
	t.Logf("Peer system cab requests after transitioning Unknown to Pending: %v", peerSystem.States[peerSystem.OwnId].CabRequests)
	orderProtocol.TransitionForCabRequestsByType(&peerSystem, transitions, orderProtocol.UnknownToAssigned)
	t.Logf("Peer system cab requests after transitioning Unknown to Assigned: %v", peerSystem.States[peerSystem.OwnId].CabRequests)
	orderProtocol.TransitionForCabRequestsByType(&peerSystem, transitions, orderProtocol.UnknownToNoOrder)
	t.Logf("Peer system cab requests after transitioning Unknown to NoOrder: %v", peerSystem.States[peerSystem.OwnId].CabRequests)

}
