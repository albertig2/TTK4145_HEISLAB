package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronisation"

	//"fmt"

	//"fmt"
	"time"
)

func orderRutine(peerView *elevatorConfig.ElevatorSystem, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, orderChannels elevatorConfig.OrderChannels, newHallOrders []elevatorConfig.ButtonEvent, newCabOrders []elevatorConfig.ButtonEvent, servicedHallOrders []elevatorConfig.ButtonEvent, servicedCabOrders []elevatorConfig.ButtonEvent) {
	HallRequestTransitions := getAllHallRequestTransitions(peerView, hallRequestsForAllElevators, newHallOrders, servicedHallOrders)
	CabRequestTransitions := getAllCabRequestTransitions(peerView, cabRequestsForAllElevators, newCabOrders, servicedCabOrders)
	transitioningAllHallRequests(peerView, HallRequestTransitions, orderChannels)
	transitioningAllCabRequests(peerView, CabRequestTransitions, orderChannels)
}

// test1
func ManagingOrders(ownId string, orderChannels elevatorConfig.OrderChannels, synchronizationChannels elevatorConfig.SynchronisationChannels, hardwareChannels elevatorConfig.ElevatorHardwareChannelsStruckt) {
	LocalPeerView := elevatorConfig.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&LocalPeerView, ownId)

	hallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus)
	cabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus)
	hallRequestsForAllElevators[ownId] = LocalPeerView.HallRequests
	cabRequestsForAllElevators[ownId] = LocalPeerView.States[ownId].CabRequests

	servicedHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	servicedCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	newHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	newCabOrders := make([]elevatorConfig.ButtonEvent, 0)

	reinitialize := false
	ticker := time.NewTicker(1 * time.Second / 2)
	defer ticker.Stop()
	for {
		select {
		case servicedorder := <-orderChannels.ServicedOrderChannel:
			if servicedorder.Button == elevatorConfig.HallUp || servicedorder.Button == elevatorConfig.HallDown {
				servicedHallOrders = append(servicedHallOrders, servicedorder)
			} else if servicedorder.Button == elevatorConfig.Cab {
				servicedCabOrders = append(servicedCabOrders, servicedorder)
			}
		case newOrder := <-orderChannels.NewRecievedOrderChannel:
			if newOrder.Button == elevatorConfig.HallUp || newOrder.Button == elevatorConfig.HallDown {
				newHallOrders = append(newHallOrders, newOrder)
			} else if newOrder.Button == elevatorConfig.Cab {
				newCabOrders = append(newCabOrders, newOrder)
			}
		case elevatorUpdate := <-synchronizationChannels.UpdateElevatorSystemWithElevatorChannel:
			synchronisation.UpdateElevatorSystemFromElevator(elevatorUpdate, &LocalPeerView)
		case externalPeerView := <-synchronizationChannels.UpdateElevatorSystemWithPeerChannel:
			synchronisation.UpdateElevatorSystemWithPeer(&LocalPeerView, &externalPeerView, hallRequestsForAllElevators, cabRequestsForAllElevators)
		case alivePeers := <-synchronizationChannels.AlivePeersChannel:
			filteredAlivePeers := []string{}
			for _, peerID := range alivePeers {
				if _, ok := LocalPeerView.States[peerID]; ok {
					filteredAlivePeers = append(filteredAlivePeers, peerID)
				}
			}
			if !synchronisation.Contains(filteredAlivePeers, LocalPeerView.OwnId) {
				filteredAlivePeers = append(filteredAlivePeers, LocalPeerView.OwnId)
			}
			synchronisation.SetAlivePeers(&LocalPeerView, filteredAlivePeers)
		case restart := <-hardwareChannels.RestartElevatorChannel:
			reinitialize = restart
			if restart {
				synchronisation.InitializeElevatorSystem(&LocalPeerView, ownId)
				hallRequestsForAllElevators[ownId] = LocalPeerView.HallRequests
				cabRequestsForAllElevators[ownId] = LocalPeerView.States[ownId].CabRequests
				servicedHallOrders = servicedHallOrders[:0]
				servicedCabOrders = servicedCabOrders[:0]
				newHallOrders = newHallOrders[:0]
				newCabOrders = newCabOrders[:0]
			}
		case <-ticker.C:
			if !reinitialize && len(LocalPeerView.AlivePeers) != 0 {
				allFloorsValid := true
				for _, peerID := range LocalPeerView.AlivePeers {
					if LocalPeerView.States[peerID].Floor == -1 {
						//fmt.Printf("Skipping ordering: elevator %s has invalid floor (-1)\n", peerID)
						allFloorsValid = false
						break
					}
				}
				if allFloorsValid {
					orderRutine(&LocalPeerView, hallRequestsForAllElevators, cabRequestsForAllElevators, orderChannels, newHallOrders, newCabOrders, servicedHallOrders, servicedCabOrders)
					hallRequestsForAllElevators[LocalPeerView.OwnId] = LocalPeerView.HallRequests
					cabRequestsForAllElevators[LocalPeerView.OwnId] = LocalPeerView.States[LocalPeerView.OwnId].CabRequests
					servicedHallOrders = servicedHallOrders[:0]
					servicedCabOrders = servicedCabOrders[:0]
					newHallOrders = newHallOrders[:0]
					newCabOrders = newCabOrders[:0]
				}
			}
		}
		select {
		case synchronizationChannels.UpdateElevatorSystemWithElevatorSystemChannel <- *synchronisation.CopyElevatorSystem(&LocalPeerView):
		default:
		}
	}
}
