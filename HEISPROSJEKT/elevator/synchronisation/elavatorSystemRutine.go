package synchronisation

import (
	"HEISPROSJEKT/elevatorConfig"
)

func elevatorSystemRutine(system *ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]OrderStatus, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]OrderStatus, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt, receivedWorldview chan ElevatorSystem) {
	for {
		select {
		case newRecievedOrder := <-elevatorOrderChannels.NewRecievedOrderChannel:
			if newRecievedOrder.Button == elevatorConfig.Cab {
				SetCabRequests(system, newRecievedOrder.Floor, Pending)
				CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests // Always update these before checking for transitions or doing anything based on the status of the floors
			} else if newRecievedOrder.Button == elevatorConfig.HallUp {
				SetHallRequests(system, newRecievedOrder.Floor, int(elevatorConfig.HallUp), Pending)
				HallRequestsForAllElevators[system.OwnId] = system.HallRequests // Always update these before checking for transitions or doing anything based on the status of the floors
			} else if newRecievedOrder.Button == elevatorConfig.HallDown {
				SetHallRequests(system, newRecievedOrder.Floor, int(elevatorConfig.HallDown), Pending)
				HallRequestsForAllElevators[system.OwnId] = system.HallRequests // Always update these before checking for transitions or doing anything based on the status of the floors
			}
		case peerSystem := <-receivedWorldview:
			UpdateElevatorSystemWithPeer(system, &peerSystem, HallRequestsForAllElevators, CabRequestsForAllElevators)
			// Add more cases for other channels if needed
		}
	}
}
