package synchronization

import (
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"time"
)

func SynchronizeElevators(elevator chan elevatorConfig.Elevator, synchronizationChannels elevatorConfig.SynchronizationChannels, controllerChannels elevatorConfig.ControllerChannels, ownId string) {
	peerView := elevatorConfig.PeerView{}
	InitializePeerView(&peerView, ownId)

	broadcastTicker := time.NewTicker(time.Second / 30)
	printTicker := time.NewTicker(time.Second * 2)
	for {
		select {
		case incommingBroadcast := <-synchronizationChannels.BcastIncomingMessagesChannel:
			if incommingBroadcast.OwnId != ownId {
				synchronizationChannels.UpdateElevatorSystemWithPeerChannel <- incommingBroadcast
			}

		case peerUpdate := <-synchronizationChannels.PeerUpdateChannel:
			debuggingHelpers.PrintPeerUpdate(peerUpdate)
			synchronizationChannels.AlivePeersChannel <- peerUpdate.Peers

		case elevatorUpdate := <-controllerChannels.LocalElevatorChannel:
			synchronizationChannels.UpdateElevatorSystemWithElevatorChannel <- elevatorUpdate

		case systemUpdate := <-synchronizationChannels.UpdateElevatorSystemWithElevatorSystemChannel:
			peerView = systemUpdate

		case <-broadcastTicker.C:
			synchronizationChannels.BcastOutgoingMessagesChannel <- peerView
			broadcastTicker.Reset(time.Second / 30)

		case <-printTicker.C:
			debuggingHelpers.PrintPeerViewUpdate(peerView)
		}
	}
}
