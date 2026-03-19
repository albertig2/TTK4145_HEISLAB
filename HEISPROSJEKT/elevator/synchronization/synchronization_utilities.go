package synchronization

/*
Provides channel initialization for synchronization
*/

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"Network-go/network/peers"
)

func InitializeSynchronizationChannels() elevatorConfig.SynchronizationChannels {
	syncronizationChannels := elevatorConfig.SynchronizationChannels{
		PeerUpdateChannel:                                  make(chan peers.PeerUpdate),
		PeerTransmitEnableChannel:                          make(chan bool),
		BroadcastIncomingMessagesChannel:                   make(chan elevatorConfig.PeerView),
		BroadcastOutgoingMessagesChannel:                   make(chan elevatorConfig.PeerView),
		UpdatePeerViewWithLocalElevatorChannel:             make(chan elevatorConfig.Elevator),
		LocalElevatorChannel:                               make(chan elevatorConfig.Elevator),
		UpdatePeerViewforBroadcastWithLocalPeerViewChannel: make(chan elevatorConfig.PeerView),
		UpdateLocalPeerViewWithExternalPeerViewChannel:     make(chan elevatorConfig.PeerView),
		AlivePeersChannel:                                  make(chan []string),
	}

	return syncronizationChannels
}
