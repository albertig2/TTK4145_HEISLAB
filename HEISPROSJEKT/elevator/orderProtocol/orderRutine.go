package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronisation"
)

func orderRutine(system *synchronisation.ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt, receivedWorldview chan synchronisation.ElevatorSystem, alivePeers []string) {
	HallRequestTransitions := GetAllHallRequestTransitions(system, HallRequestsForAllElevators, alivePeers)
	CabRequestTransitions := GetAllCabRequestTransitions(system, CabRequestsForAllElevators, alivePeers)
	TransitioningAllHallRequests(system, HallRequestTransitions, alivePeers, elevatorOrderChannels)
	TransitioningAllCabRequests(system, CabRequestTransitions, elevatorOrderChannels)
	HallRequestsForAllElevators[system.OwnId] = system.HallRequests
	CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests
}
