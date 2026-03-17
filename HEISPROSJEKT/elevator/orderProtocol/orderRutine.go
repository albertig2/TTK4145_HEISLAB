package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronization"

	//"fmt"

	//"fmt"
	"time"
)

func orderRutine(system *elevatorConfig.PeerView,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus,
	elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt,
	newHallOrders []elevatorConfig.ButtonEvent,
	newCabOrders []elevatorConfig.ButtonEvent,
	servicedHallOrders []elevatorConfig.ButtonEvent,
	servicedCabOrders []elevatorConfig.ButtonEvent) {

	HallRequestTransitions := GetAllHallRequestTransitions(system, HallRequestsForAllElevators, newHallOrders, servicedHallOrders)

	CabRequestTransitions := GetAllCabRequestTransitions(system, CabRequestsForAllElevators, newCabOrders, servicedCabOrders)

	TransitioningAllHallRequests(system, HallRequestTransitions, elevatorOrderChannels)

	TransitioningAllCabRequests(system, CabRequestTransitions, elevatorOrderChannels)

}

// test1
func RunOrder(
	id string,
	elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt,
	synchronisationChannels elevatorConfig.SynchronizationChannels,
	hardwareChannel elevatorConfig.ElevatorHardwareChannelsStruckt,
) {

	system := elevatorConfig.PeerView{}
	synchronization.InitializePeerView(&system, id)

	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus)
	HallRequestsForAllElevators[id] = system.HallRequests
	CabRequestsForAllElevators[id] = system.States[id].CabRequests

	servicedHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	servicedCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	newHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	newCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	UpdatedSystem := elevatorConfig.PeerView{}
	synchronization.InitializePeerView(&UpdatedSystem, id)

	paused := false
	ticker := time.NewTicker(1 * time.Second / 2)
	//ticker.Stop()
	defer ticker.Stop()
	for {
		//fmt.Println("Before select in order routine")
		select {
		case servicedorder := <-elevatorOrderChannels.ServicedOrderChannel:
			//fmt.Println("Before servidec order")
			if servicedorder.Button == elevatorConfig.HallUp || servicedorder.Button == elevatorConfig.HallDown {
				servicedHallOrders = append(servicedHallOrders, servicedorder)
			} else if servicedorder.Button == elevatorConfig.Cab {
				servicedCabOrders = append(servicedCabOrders, servicedorder)
			}
			//fmt.Println("After serviced order")
		case newOrder := <-elevatorOrderChannels.NewRecievedOrderChannel:
			//fmt.Println("Before new order")
			if newOrder.Button == elevatorConfig.HallUp || newOrder.Button == elevatorConfig.HallDown {
				newHallOrders = append(newHallOrders, newOrder)
			} else if newOrder.Button == elevatorConfig.Cab {
				newCabOrders = append(newCabOrders, newOrder)
			}
			//fmt.Println("After new order")
		// etterhvert kan denne fjernes
		case elevatorUpdate := <-synchronisationChannels.UpdateElevatorSystemWithElevatorChannel:
			synchronization.UpdateElevatorSystemFromElevator(elevatorUpdate, &system)
		case PeerSystem := <-synchronisationChannels.UpdateElevatorSystemWithPeerChannel:
			synchronization.UpdateLocalPeerViewWithPeer(&system, &PeerSystem, HallRequestsForAllElevators, CabRequestsForAllElevators)
		case AlivePeers := <-synchronisationChannels.AlivePeersChannel:
			filteredAlivePeers := []string{}
			for _, peerID := range AlivePeers {
				if _, ok := system.States[peerID]; ok {
					filteredAlivePeers = append(filteredAlivePeers, peerID)
				}
			}

			if !synchronization.Contains(filteredAlivePeers, system.OwnId) {
				filteredAlivePeers = append(filteredAlivePeers, system.OwnId)
			}
			synchronization.SetAlivePeers(&system, filteredAlivePeers)
		case motorstop := <-hardwareChannel.RestartElevatorChannel:
			//fmt.Printf("Received motor failure status: %v", motorstop)
			paused = motorstop
			if motorstop {
				synchronization.InitializePeerView(&system, id)
				HallRequestsForAllElevators[id] = system.HallRequests
				CabRequestsForAllElevators[id] = system.States[id].CabRequests
				servicedHallOrders = servicedHallOrders[:0]
				servicedCabOrders = servicedCabOrders[:0]
				newHallOrders = newHallOrders[:0]
				newCabOrders = newCabOrders[:0]
			}
			//fmt.Printf("Paused status: %v", paused)
		case <-ticker.C:

			if !paused && len(system.AlivePeers) != 0 {
				// Check that all elevators have a valid floor before ordering
				allFloorsValid := true
				for _, peerID := range system.AlivePeers {
					if system.States[peerID].Floor == -1 {
						//fmt.Printf("Skipping ordering: elevator %s has invalid floor (-1)\n", peerID)
						allFloorsValid = false
						break
					}
				}
				if allFloorsValid {
					//fmt.Println("Before ordering")
					orderRutine(&system, HallRequestsForAllElevators, CabRequestsForAllElevators, elevatorOrderChannels, newHallOrders, newCabOrders, servicedHallOrders, servicedCabOrders)
					HallRequestsForAllElevators[system.OwnId] = system.HallRequests
					CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests
					servicedHallOrders = servicedHallOrders[:0]
					servicedCabOrders = servicedCabOrders[:0]
					newHallOrders = newHallOrders[:0]
					newCabOrders = newCabOrders[:0]
					//fmt.Println("After ordering")
				}
			}
		}
		//fmt.Println("After select in orderrutine")
		select {
		case synchronisationChannels.UpdateElevatorSystemWithElevatorSystemChannel <- *synchronization.CopyPeerView(&system):
			//fmt.Println("Sent system update to synchroniseElevators (deep copy)")
		default:
			//fmt.Println("Channel full, system update dropped")
		}
	}
}

// Race conditions for HallRequests matriser og cabrequests og ElevatorSystem struct.....
