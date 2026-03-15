package synchronisation

import (
	// "HEISPROSJEKT/debuggingHelpers"
	// "HEISPROSJEKT/elevatorConfig"
	// "Network-go/network/peers"

	// "fmt"
	// "strconv"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"

	"fmt"

	//"fmt"
	"time"
)

// called as a go functions
func SynchroniseElevators(elevatorUpdates chan elevatorConfig.Elevator, synchronisationChannels elevatorConfig.SynchronisationChannels, ownId string) {

	elevatorSystem := elevatorConfig.ElevatorSystem{}
	InitializeElevatorSystem(&elevatorSystem, ownId)

	broadcastTicker := time.NewTicker(time.Second / 30)

	for {
		fmt.Println("Before select in synchroniseElevators")

		select {
		case incommingBroadcast := <-synchronisationChannels.BcastIncomingMessagesChannel:
			fmt.Println("Received broadcast in synchroniseElevators")
			if incommingBroadcast.OwnId != ownId {
				synchronisationChannels.UpdateElevatorSystemWithPeerChannel <- incommingBroadcast
			}

		case peerUpdate := <-synchronisationChannels.PeerUpdateChannel:
			fmt.Println("Recieved peer list")
			synchronisationChannels.AlivePeersChannel <- peerUpdate.Peers
			fmt.Println("Sent alive peers to orders")
			debuggingHelpers.PrintPeerUpdate(peerUpdate)

		case elevatorUpdate := <-elevatorUpdates:
			fmt.Println("Received elevator struct from elevator fsm")
			synchronisationChannels.UpdateElevatorSystemWithElevatorChannel <- elevatorUpdate
			fmt.Println("Sent elevator struct to orders")
			// Maybe I should use this too in orders(?)

		case systemUpdate := <-synchronisationChannels.UpdateElevatorSystemWithElevatorSystemChannel:
			elevatorSystem = systemUpdate
			fmt.Println("Received system update in synchroniseElevators")

		case <-broadcastTicker.C:
			fmt.Println("Sending broadcast in synchroniseElevators")
			synchronisationChannels.BcastOutgoingMessagesChannel <- elevatorSystem
			fmt.Println("Sent broadcast in synchroniseElevators")
			broadcastTicker.Reset(time.Second / 30)
		}
		fmt.Println("After select in synchroniseElevators")

	}

	//run all synchronisation on events (like elevator fsm, run events on most sync channel?)
	//recieve eleavtor from FSM
	//recieve update from peerNetwork
	//recieve broadcast form peer

}
