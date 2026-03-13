package synchronisation

import (
	// "HEISPROSJEKT/debuggingHelpers"
	// "HEISPROSJEKT/elevatorConfig"
	// "Network-go/network/peers"

	// "fmt"
	// "strconv"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"

	//"fmt"
	"time"
)

// called as a go functions
func SynchroniseElevators(elevatorUpdates chan elevatorConfig.Elevator, synchronisationChannels elevatorConfig.SynchronisationChannels, ownId string) {

	elevatorSystem := elevatorConfig.ElevatorSystem{}
	InitializeElevatorSystem(&elevatorSystem, ownId)

	broadcastTicker := time.NewTicker(time.Second / 30)

	//aliveList := []string{}
	hallRequestsForAllElevators := map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{}
	cabRequestsForAllElevators := map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{}

	for {
		select {

		case incommingBroadcast := <-synchronisationChannels.BcastIncomingMessagesChannel:
			if incommingBroadcast.OwnId == ownId {
				continue
			}
			UpdateElevatorSystemWithPeer(&elevatorSystem, &incommingBroadcast, hallRequestsForAllElevators, cabRequestsForAllElevators)
			peerRequests := elevatorConfig.PeerRequestUpdate{
				PeerID:   incommingBroadcast.OwnId,
				HallReqs: hallRequestsForAllElevators[incommingBroadcast.OwnId],
				CabReqs:  cabRequestsForAllElevators[incommingBroadcast.OwnId],
			}
			synchronisationChannels.UpdatePeerRequests <- peerRequests

		case peerUpdate := <-synchronisationChannels.PeerUpdateChl:
			filteredAlivePeers := []string{}
			for _, peerID := range peerUpdate.Peers {
				if _, ok := elevatorSystem.States[peerID]; ok {
					filteredAlivePeers = append(filteredAlivePeers, peerID)
				}
			}
			SetAlivePeers(&elevatorSystem, filteredAlivePeers)
			debuggingHelpers.PrintPeerUpdate(peerUpdate)

		case elevatorUpdate := <-elevatorUpdates:
			// Maybe I should use this too in orders(?)
			UpdateElevatorSystemFromELevator(elevatorUpdate, &elevatorSystem)

		case systemUpdate := <-synchronisationChannels.UpdateElevatorSystem:
			elevatorSystem = systemUpdate

		case <-broadcastTicker.C:

			synchronisationChannels.BcastOutgoingMessagesChannel <- elevatorSystem
			broadcastTicker.Reset(time.Second / 30)
		}

		select {
		case synchronisationChannels.UpdateElevatorSystem <- elevatorSystem:
		default:
		}
	}

	//run all synchronisation on events (like elevator fsm, run events on most sync channel?)
	//recieve eleavtor from FSM
	//recieve update from peerNetwork
	//recieve broadcast form peer

}
