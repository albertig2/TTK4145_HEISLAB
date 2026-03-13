package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronisation"

	//"fmt"
	"time"
)

func orderRutine(system *elevatorConfig.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus,
	elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt,
	newHallOrders []elevatorConfig.ButtonEvent,
	newCabOrders []elevatorConfig.ButtonEvent,
	servicedHallOrder elevatorConfig.ButtonEvent,
	servicedCabOrder elevatorConfig.ButtonEvent) {

	HallRequestTransitions := GetAllHallRequestTransitions(system, HallRequestsForAllElevators, newHallOrders, servicedHallOrder)

	CabRequestTransitions := GetAllCabRequestTransitions(system, CabRequestsForAllElevators, newCabOrders, servicedCabOrder)

	TransitioningAllHallRequests(system, HallRequestTransitions, elevatorOrderChannels)

	TransitioningAllCabRequests(system, CabRequestTransitions, elevatorOrderChannels)

}

// test1
func RunOrder(
	id string,
	elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt,
	synchronisationChannels elevatorConfig.SynchronisationChannels,
	hardwareChannel elevatorConfig.ElevatorHardwareChannelsStruckt,
) {

	system := elevatorConfig.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&system, id)

	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus)
	HallRequestsForAllElevators[id] = system.HallRequests
	CabRequestsForAllElevators[id] = system.States[id].CabRequests

	servicedHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	servicedCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	newHallOrders := make([]elevatorConfig.ButtonEvent, 0)
	newCabOrders := make([]elevatorConfig.ButtonEvent, 0)
	UpdatedSystem := elevatorConfig.ElevatorSystem{}
	synchronisation.InitializeElevatorSystem(&UpdatedSystem, id)

	paused := false
	ticker := time.NewTicker(1 * time.Second)
	//ticker.Stop()
	defer ticker.Stop()
	for {
		select {
		case servicedorder := <-elevatorOrderChannels.ServicedOrderChannel:

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
		// etterhvert kan denne fjernes
		case UpdatedSystem = <-synchronisationChannels.UpdateElevatorSystem:
			system = UpdatedSystem
		case UpdatedPeerRequests := <-synchronisationChannels.UpdatePeerRequests:
			HallRequestsForAllElevators[UpdatedPeerRequests.PeerID] = UpdatedPeerRequests.HallReqs
			CabRequestsForAllElevators[UpdatedPeerRequests.PeerID] = UpdatedPeerRequests.CabReqs
		case motorstop := <-hardwareChannel.MotorFailureChannel:
			paused = motorstop
			if motorstop {
				synchronisation.InitializeElevatorSystem(&system, id)
				HallRequestsForAllElevators[id] = system.HallRequests
				CabRequestsForAllElevators[id] = system.States[id].CabRequests
			}
		case <-ticker.C:

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

				orderRutine(&system, HallRequestsForAllElevators, CabRequestsForAllElevators, elevatorOrderChannels, newHallOrders, newCabOrders, sh, sc)
				HallRequestsForAllElevators[system.OwnId] = system.HallRequests
				CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests
				servicedHallOrders = servicedHallOrders[:0]
				servicedCabOrders = servicedCabOrders[:0]
				newHallOrders = newHallOrders[:0]
				newCabOrders = newCabOrders[:0]

				synchronisationChannels.UpdateElevatorSystem <- system
			}

		}
	}
}

// Race conditions for HallRequests matriser og cabrequests og ElevatorSystem struct.....
