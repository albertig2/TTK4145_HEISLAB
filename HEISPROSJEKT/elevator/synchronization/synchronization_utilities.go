package synchronization

/*
Provides channel initialization for synchronization
*/

import (
	"HEISPROSJEKT/elevatorConfig"
	"Network-go/network/peers"
)

func InitializeSynchronizationChannels() elevatorConfig.SynchronizationChannels {
	syncronizationChannels := elevatorConfig.SynchronizationChannels{
		PeerUpdateChannel:                             make(chan peers.PeerUpdate),
		PeerTxEnableChannel:                           make(chan bool),
		BcastIncomingMessagesChannel:                  make(chan elevatorConfig.PeerView),
		BcastOutgoingMessagesChannel:                  make(chan elevatorConfig.PeerView),
		UpdateElevatorSystemWithElevatorChannel:       make(chan elevatorConfig.Elevator),
		UpdateElevatorSystemWithElevatorSystemChannel: make(chan elevatorConfig.PeerView),
		UpdateElevatorSystemWithPeerChannel:           make(chan elevatorConfig.PeerView),
		AlivePeersChannel:                             make(chan []string),
	}

	return syncronizationChannels
}
