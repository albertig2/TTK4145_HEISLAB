package synchronization

import (
	"HEISPROSJEKT/elevatorConfig"
	"Network-go/network/peers"
)

var (
	alivePeersList []string
	deadPeersList  []string
)

func InitializeSynchrinizationChannels() elevatorConfig.SynchronizationChannels {
	channels := elevatorConfig.SynchronizationChannels{
		PeerUpdateChannel:                             make(chan peers.PeerUpdate),
		PeerTxEnableChannel:                           make(chan bool),
		BcastIncomingMessagesChannel:                  make(chan elevatorConfig.PeerView),
		BcastOutgoingMessagesChannel:                  make(chan elevatorConfig.PeerView),
		UpdateElevatorSystemWithElevatorChannel:       make(chan elevatorConfig.Elevator),
		UpdateElevatorSystemWithElevatorSystemChannel: make(chan elevatorConfig.PeerView),
		UpdateElevatorSystemWithPeerChannel:           make(chan elevatorConfig.PeerView),
		AlivePeersChannel:                             make(chan []string),
	}

	return channels
}

func UpdatePeerList(synchronizationChannelss elevatorConfig.SynchronizationChannels) {
	for {
		peerUpdate := <-synchronizationChannelss.PeerUpdateChannel

		alivePeersList = peerUpdate.Peers

		deadPeersList = append(deadPeersList, peerUpdate.Lost...)

		newPeer := peerUpdate.New

		if newPeer != "" {
			newList := []string{}

			for _, dead := range deadPeersList {
				if dead != newPeer {
					newList = append(newList, dead)
				}
			}

			deadPeersList = newList
		}
	}
}



func UpdateElevatorSystemFromElevator(elevator elevatorConfig.Elevator, peerView *elevatorConfig.PeerView) {
	SetBehavior(peerView, elevatorConfig.Behavior(elevator.Behavior))
	SetDirection(peerView, elevator.Direction)
	//fmt.Printf("Updating elevator system from elevator struct. Elevator floor: %d\n", elevator.Floor)
	if elevator.Floor >= 0 && elevator.Floor < elevatorConfig.N_FLOORS {
		SetFloor(peerView, elevator.Floor)
	}
}