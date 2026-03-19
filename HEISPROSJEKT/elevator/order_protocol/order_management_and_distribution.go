package orderProtocol

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"HEISPROSJEKT/synchronization"
	"time"
)

/*
This file contains the main routine for managing and distributing orders in the elevator system.
All communication regarding orders, including new orders, serviced orders, and order state updates, is handled here.
The main function, ManageAndDistributeOrders, coordinates order processing, synchronization protocol interaction,
and communication with the elevator controller. It updates the local peer's view with the state of the local elevator and
the external elevators.
*/

func ManageAndDistributeOrders(ownId string, orderChannels elevatorConfig.OrderChannels, synchronizationChannels elevatorConfig.SynchronizationChannels, hardwareChannels elevatorConfig.ControllerChannels) {
	localPeerView := elevatorConfig.PeerView{}
	hallOrdersForAllElevators, cabOrdersForAllElevators := initializePeerView(&localPeerView, ownId)
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
			synchronization.UpdateLocalPeerViewWithPeer(&localPeerView, &externalPeerView, hallOrdersForAllElevators, cabOrdersForAllElevators)

		case alivePeers := <-synchronizationChannels.AlivePeersChannel:
			validAlivePeers := filterValidPeersAndIncludeOwnId(&localPeerView, alivePeers)
			synchronization.SetAlivePeers(&localPeerView, validAlivePeers)

		case restart := <-hardwareChannels.RestartElevatorChannel:
			reinitialize = restart
			if restart {
				hallOrdersForAllElevators, cabOrdersForAllElevators = initializePeerView(&localPeerView, ownId)
				servicedHallOrders, servicedCabOrders, newHallOrders, newCabOrders = initializeOrDrainOrders()
			}

		case <-ticker.C:
			if !reinitialize && areAllFloorsValid(&localPeerView) {
				orderRutine(&localPeerView, &hallOrdersForAllElevators, &cabOrdersForAllElevators, orderChannels, newHallOrders, newCabOrders, servicedHallOrders, servicedCabOrders)
				servicedHallOrders, servicedCabOrders, newHallOrders, newCabOrders = initializeOrDrainOrders()
			}
			
		}
		select {
		case synchronizationChannels.UpdatePeerViewforBroadcastWithLocalPeerViewChannel <- *synchronization.CopyPeerView(&localPeerView):
		default:
		}
	}
}
