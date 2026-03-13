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
func SynchroniseElevators(elevatorUpdates chan elevatorConfig.Elevator, synchronisationChannels SynchronisationChannels, ownId string) {

	elevatorSystem := ElevatorSystem{}
	InitializeElevatorSystem(&elevatorSystem, ownId)

	broadcastTicker := time.NewTicker(time.Second / 30)

	//aliveList := []string{}
	hallRequestsForAllElevators := map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{}
	cabRequestsForAllElevators := map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{}

	for {
		select {

		case incommingBroadcast := <-synchronisationChannels.BcastIncomingMessagesChannel:

			UpdateElevatorSystemWithPeer(&elevatorSystem, &incommingBroadcast, hallRequestsForAllElevators, cabRequestsForAllElevators)

		case peerUpdate := <-synchronisationChannels.PeerUpdateChl:

			//aliveList = peerUpdate.Peers
			debuggingHelpers.PrintPeerUpdate(peerUpdate)

		case elevatorUpdate := <-elevatorUpdates:

			UpdateElevatorSystemFromELevator(elevatorUpdate, &elevatorSystem)
		
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
