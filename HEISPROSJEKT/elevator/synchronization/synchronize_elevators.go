package synchronization

/*
Runs the main synchronization loop.
Coordinates communication between peers, local elevator updates,
and periodic broadcasting of the system state.
*/

import (
	"HEISPROSJEKT/debuggingHelpers"
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"time"
)

func SynchronizeElevators(synchronizationChannels elevatorConfig.SynchronizationChannels, controllerChannels elevatorConfig.ControllerChannels, ownId string) {
	peerViewForBroadcast := elevatorConfig.PeerView{}
	InitializePeerView(&peerViewForBroadcast, ownId)

	broadcastTicker := time.NewTicker(time.Second / 30)
	printTicker := time.NewTicker(time.Second * 2)
	for {
		select {
		case externalPeerView := <-synchronizationChannels.BroadcastIncomingMessagesChannel:
			if externalPeerView.OwnId != ownId {
				synchronizationChannels.UpdateLocalPeerViewWithExternalPeerViewChannel <- externalPeerView
			}

		case peerUpdate := <-synchronizationChannels.PeerUpdateChannel:
			debuggingHelpers.PrintPeerUpdate(peerUpdate)
			synchronizationChannels.AlivePeersChannel <- peerUpdate.Peers

		case elevatorUpdate := <-synchronizationChannels.LocalElevatorChannel:
			synchronizationChannels.UpdatePeerViewWithLocalElevatorChannel <- elevatorUpdate

		case localPeerView := <-synchronizationChannels.UpdatePeerViewforBroadcastWithLocalPeerViewChannel:
			peerViewForBroadcast = localPeerView

		case <-broadcastTicker.C:
			synchronizationChannels.BroadcastOutgoingMessagesChannel <- peerViewForBroadcast
			broadcastTicker.Reset(time.Second / 30)

		case <-printTicker.C:
			debuggingHelpers.PrintPeerViewUpdate(peerViewForBroadcast)
		}
	}
}
