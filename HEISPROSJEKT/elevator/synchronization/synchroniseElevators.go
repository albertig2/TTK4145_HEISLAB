package synchronization

import (
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"time"
)


func SynchronizeElevators(elevator chan elevatorConfig.Elevator, synchronizationChannels elevatorConfig.SynchronizationChannels, ownId string) {
	elevatorSystem := elevatorConfig.PeerView{}
	InitializeElevatorSystem(&elevatorSystem, ownId)

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

		case elevatorUpdate := <-elevator:
			synchronizationChannels.UpdateElevatorSystemWithElevatorChannel <- elevatorUpdate

		case systemUpdate := <-synchronizationChannels.UpdateElevatorSystemWithElevatorSystemChannel:
			elevatorSystem = systemUpdate

		case <-broadcastTicker.C:
			synchronizationChannels.BcastOutgoingMessagesChannel <- elevatorSystem
			broadcastTicker.Reset(time.Second / 30)

		case <-printTicker.C:
			debuggingHelpers.PrintElevatorSystem(elevatorSystem)
		}
	}
}
