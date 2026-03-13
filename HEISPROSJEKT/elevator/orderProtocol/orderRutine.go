package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronisation"
	"fmt"
	"time"
)

func orderRutine(system *synchronisation.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus,
	elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt,
	newHallOrders []elevatorConfig.ButtonEvent,
	newCabOrders []elevatorConfig.ButtonEvent,
	servicedHallOrder elevatorConfig.ButtonEvent,
	servicedCabOrder elevatorConfig.ButtonEvent,
	alivePeers []string) {
	fmt.Println("Start hallrequestTransitions")
	HallRequestTransitions := GetAllHallRequestTransitions(system, HallRequestsForAllElevators, newHallOrders, servicedHallOrder, alivePeers)
	fmt.Println("Start CabrequestTransitions")
	CabRequestTransitions := GetAllCabRequestTransitions(system, CabRequestsForAllElevators, newCabOrders, servicedCabOrder, alivePeers)
	fmt.Println("Start Transition all hall requests")
	TransitioningAllHallRequests(system, HallRequestTransitions, alivePeers, elevatorOrderChannels)
	fmt.Println("Start Transition all Cab requests")
	TransitioningAllCabRequests(system, CabRequestTransitions, elevatorOrderChannels)
	fmt.Println("End Transition all Cab requests")
}

func RunOrder(
	id string,
	elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt,
	synchronisationChannels synchronisation.SynchronisationChannels,
	hardwareChannel elevatorConfig.ElevatorHardwareChannelsStruckt,
) {

	system := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system, id)

	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus)
	HallRequestsForAllElevators[id] = system.HallRequests
	CabRequestsForAllElevators[id] = system.States[id].CabRequests

	servicedHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	servicedCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	newHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	newCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	alivePeers := make([]string, 0)
	UpdatedSystem := synchronisation.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&UpdatedSystem, id)
	alivePeers = []string{id}

	paused := false
	ticker := time.NewTicker(1 * time.Second)
	//ticker.Stop()
	defer ticker.Stop()
	for {
		select {
		case servicedorder := <-elevatorOrderChannels.ServicedOrderChannel:
			fmt.Println("Recieved order", servicedorder, "In orderroutine")
			if servicedorder.Button == elevatorConfig.HallUp || servicedorder.Button == elevatorConfig.HallDown {
				servicedHallOrders = append(servicedHallOrders, servicedorder)
			} else if servicedorder.Button == elevatorConfig.Cab {
				servicedCabOrders = append(servicedCabOrders, servicedorder)
			}
		case newOrder := <-elevatorOrderChannels.NewRecievedOrderChannel:
			if newOrder.Button == elevatorConfig.HallUp || newOrder.Button == elevatorConfig.HallDown {
				newHallOrders = append(newHallOrders, newOrder)
			} else if newOrder.Button == elevatorConfig.Cab {
				newCabOrders = append(newCabOrders, newOrder)
			}
		case peers := <-synchronisationChannels.PeerUpdateChl:
			alivePeers = peers.Peers
		case peerSystem := <-synchronisationChannels.BcastIncomingMessagesChannel:
			synchronisation.UpdateElevatorSystemWithPeer(&system, &peerSystem, HallRequestsForAllElevators, CabRequestsForAllElevators)
		case UpdatedSystem = <-synchronisationChannels.UpdateElevatorSystem:
			system = UpdatedSystem
		case motorstop := <-hardwareChannel.MotorFailureChannel:
			paused = motorstop
			if motorstop {
				synchronisation.InitializeElevatorSystem(&system, id)
				HallRequestsForAllElevators[id] = system.HallRequests
				CabRequestsForAllElevators[id] = system.States[id].CabRequests
			}
		case <-ticker.C:
			fmt.Println("Ticker kicked in")
			if !paused {
				var sh, sc elevatorConfig.ButtonEvent
				sh.Floor = -1
				sc.Floor = -1
				if len(servicedHallOrders) > 0 {
					sh = servicedHallOrders[0]
				}
				if len(servicedCabOrders) > 0 {
					sc = servicedCabOrders[0]
				}
				fmt.Println("Ticker print 2")
				orderRutine(&system, HallRequestsForAllElevators, CabRequestsForAllElevators, elevatorOrderChannels, newHallOrders, newCabOrders, sh, sc, alivePeers)
				HallRequestsForAllElevators[system.OwnId] = system.HallRequests
				CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests
				servicedHallOrders = servicedHallOrders[:0]
				servicedCabOrders = servicedCabOrders[:0]
				newHallOrders = newHallOrders[:0]
				newCabOrders = newCabOrders[:0]
				
				synchronisationChannels.UpdateElevatorSystem <- system
				fmt.Println("Ticker end")
			}
			
		}
	}
}

// Race conditions for HallRequests matriser og cabrequests og ElevatorSystem struct.....
