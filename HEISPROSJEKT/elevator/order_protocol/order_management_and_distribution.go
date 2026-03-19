package orderProtocol

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"HEISPROSJEKT/synchronization"
	"time"
	//"fmt"
)

func ManageAndDistributeOrders(ownId string, orderChannels elevatorConfig.OrderChannels, synchronizationChannels elevatorConfig.SynchronizationChannels, hardwareChannels elevatorConfig.ControllerChannels) {
	localPeerView := elevatorConfig.PeerView{}
	hallRequestsForAllElevators, cabRequestsForAllElevators := initializePeerView(&localPeerView, ownId)
	servicedHallOrders, servicedCabOrders, newHallOrders, newCabOrders := initializeOrDrainOrders()
	reinitialize := false
	ticker := time.NewTicker(1 * time.Second / 2)
	defer ticker.Stop()
	for {
		select {
		case servicedorder := <-orderChannels.ServicedOrderChannel:
			servicedHallOrders, servicedCabOrders = appendOrderByType(servicedHallOrders, servicedCabOrders, servicedorder)
		case newOrder := <-orderChannels.NewRecievedOrderChannel:
			newHallOrders, newCabOrders = appendOrderByType(newHallOrders, newCabOrders, newOrder)
		case elevatorUpdate := <-synchronizationChannels.UpdatePeerViewWithLocalElevatorChannel:
			synchronization.UpdateElevatorSystemFromElevator(elevatorUpdate, &localPeerView)
		case externalPeerView := <-synchronizationChannels.UpdateLocalPeerViewWithExternalPeerViewChannel:
			synchronization.UpdateLocalPeerViewWithPeer(&localPeerView, &externalPeerView, hallRequestsForAllElevators, cabRequestsForAllElevators)
		case alivePeers := <-synchronizationChannels.AlivePeersChannel:
			validAlivePeers := filterValidPeersAndIncludeOwnId(&localPeerView, alivePeers)
			synchronization.SetAlivePeers(&localPeerView, validAlivePeers)
		case restart := <-hardwareChannels.RestartElevatorChannel:
			reinitialize = restart
			if restart {
				hallRequestsForAllElevators, cabRequestsForAllElevators = initializePeerView(&localPeerView, ownId)
				servicedHallOrders, servicedCabOrders, newHallOrders, newCabOrders = initializeOrDrainOrders()
			}
		case <-ticker.C:
			if !reinitialize && areAllFloorsValid(&localPeerView) {
				orderRutine(&localPeerView, &hallRequestsForAllElevators, &cabRequestsForAllElevators, orderChannels, newHallOrders, newCabOrders, servicedHallOrders, servicedCabOrders)
				servicedHallOrders, servicedCabOrders, newHallOrders, newCabOrders = initializeOrDrainOrders()
			}
		}
		select {
		case synchronizationChannels.UpdatePeerViewforBroadcastWithLocalPeerViewChannel <- *synchronization.CopyPeerView(&localPeerView):
		default:
		}
	}
}
