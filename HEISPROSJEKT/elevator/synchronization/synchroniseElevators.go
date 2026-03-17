package synchronization

import (
	// "HEISPROSJEKT/debuggingHelpers"
	// "HEISPROSJEKT/elevatorConfig"
	// "Network-go/network/peers"

	// "fmt"
	// "strconv"
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"

	//"fmt"

	//"fmt"
	"time"
)

// called as a go functions
func SynchroniseElevators(elevator chan elevatorConfig.Elevator, synchronizationChannels elevatorConfig.SynchronizationChannels, ownId string) {

	elevatorSystem := elevatorConfig.PeerView{}
	InitializeElevatorSystem(&elevatorSystem, ownId)

	broadcastTicker := time.NewTicker(time.Second / 30)
	printTicker := time.NewTicker(time.Second * 2)

	for {
		//fmt.Println("Before select in synchroniseElevators")

		select {
		case incommingBroadcast := <-synchronizationChannels.BcastIncomingMessagesChannel:
			//fmt.Println("Received broadcast in synchroniseElevators")
			if incommingBroadcast.OwnId != ownId {
				synchronizationChannels.UpdateElevatorSystemWithPeerChannel <- incommingBroadcast
			}

		case peerUpdate := <-synchronizationChannels.PeerUpdateChannel:
			//fmt.Println("Recieved peer list")
			debuggingHelpers.PrintPeerUpdate(peerUpdate)
			synchronizationChannels.AlivePeersChannel <- peerUpdate.Peers
			//fmt.Println("Sent alive peers to orders")

		case elevatorUpdate := <-elevator:
			//fmt.Println("Received elevator struct from elevator fsm")
			synchronizationChannels.UpdateElevatorSystemWithElevatorChannel <- elevatorUpdate
			//fmt.Println("Sent elevator struct to orders")
			// Maybe I should use this too in orders(?)

		case systemUpdate := <-synchronizationChannels.UpdateElevatorSystemWithElevatorSystemChannel:
			elevatorSystem = systemUpdate
			//fmt.Println("Received system update in synchroniseElevators")

		case <-broadcastTicker.C:
			//fmt.Println("Sending broadcast in synchroniseElevators")
			synchronizationChannels.BcastOutgoingMessagesChannel <- elevatorSystem
			//fmt.Println("Sent broadcast in synchroniseElevators")
			broadcastTicker.Reset(time.Second / 30)

		case <-printTicker.C:
			debuggingHelpers.PrintElevatorSystem(elevatorSystem)
		}
		//fmt.Println("After select in synchroniseElevators")

	}

	//run all synchronisation on events (like elevator fsm, run events on most sync channel?)
	//recieve eleavtor from FSM
	//recieve update from peerNetwork
	//recieve broadcast form peer

}
