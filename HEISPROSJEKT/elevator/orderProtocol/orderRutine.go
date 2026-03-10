package orderProtocol

import (
	"HEISPROSJEKT/communication"
	"HEISPROSJEKT/elevatorHardware"
)

func orderRutine(system *elevatorHardware.ElevatorSystem, receivedWorldView chan string) {

	peerSystem := communication.DecodeElevatorSystem(<-receivedWorldView)
	elevatorHardware.AddPeer(system, peerSystem)
	elevatorHardware.UpdateElevatorSystemFromPeer(system, &peerSystem, nil, nil)
}
