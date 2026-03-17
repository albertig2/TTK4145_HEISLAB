package synchronization

/*
func elevatorSystemRutine(id string, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt, receivedWorldview chan ElevatorSystem, hardwareChannels ElevatorHardwareChannelsStruckt) {
	system := ElevatorSystem{}
	InitializeElevatorSystem(&system, id)
	HallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus)
	CabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus)

	for {
		select {
		case peerSystem := <-receivedWorldview:
			UpdateElevatorSystemWithPeer(&system, &peerSystem, HallRequestsForAllElevators, CabRequestsForAllElevators)
			// Add more cases for other channels if needed
		case motorstop := <-hardwareChannels.MotorFailureChannel:
			if motorstop {
				InitializeElevatorSystem(&system, id)
			}
		case cabOrdersSet :<-
		}
	}
}
*/
/*
	case newRecievedOrder := <-elevatorOrderChannels.NewRecievedOrderChannel:
		if newRecievedOrder.Button == elevatorConfig.Cab {
			SetCabRequests(&system, newRecievedOrder.Floor, Pending)
			CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests // Always update these before checking for transitions or doing anything based on the status of the floors
		} else if newRecievedOrder.Button == elevatorConfig.HallUp {
			SetHallRequests(&system, newRecievedOrder.Floor, int(elevatorConfig.HallUp), Pending)
			HallRequestsForAllElevators[system.OwnId] = system.HallRequests // Always update these before checking for transitions or doing anything based on the status of the floors
		} else if newRecievedOrder.Button == elevatorConfig.HallDown {
			SetHallRequests(&system, newRecievedOrder.Floor, int(elevatorConfig.HallDown), Pending)
			HallRequestsForAllElevators[system.OwnId] = system.HallRequests // Always update these before checking for transitions or doing anything based on the status of the floors
		}
*/
// Må kanskje be elevatorSystem om å oppdatere hallrequest og cabrequests matriser?
// Fikse at serviced er noe jeg mottar ikke sender på
// Samle som en go rutine
