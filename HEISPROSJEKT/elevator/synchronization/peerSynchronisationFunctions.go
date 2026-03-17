package synchronization

import (
	"HEISPROSJEKT/elevatorConfig"
	"Network-go/network/peers"

	"fmt"
	"strconv"
	"time"
)

var (
	alivePeersList []string
	deadPeersList  []string
)

func InitializeSynchrinizationChannels() elevatorConfig.SynchronizationChannels {
	channels := elevatorConfig.SynchronizationChannels{
		PeerUpdateChannel:                             make(chan peers.PeerUpdate),
		PeerTxEnableChannel:                           make(chan bool),
		BcastIncomingMessagesChannel:                  make(chan elevatorConfig.PeerView),
		BcastOutgoingMessagesChannel:                  make(chan elevatorConfig.PeerView),
		UpdateElevatorSystemWithElevatorChannel:       make(chan elevatorConfig.Elevator),
		UpdateElevatorSystemWithElevatorSystemChannel: make(chan elevatorConfig.PeerView),
		UpdateElevatorSystemWithPeerChannel:           make(chan elevatorConfig.PeerView),
		AlivePeersChannel:                             make(chan []string),
	}

	return channels
}

// remove
func StartPeerNetworking(port int, id int, channels elevatorConfig.SynchronizationChannels) {
	fmt.Println("start")

	go peers.Receiver(port, channels.PeerUpdateChannel)
	go peers.Transmitter(port, strconv.Itoa(id), channels.PeerTxEnableChannel)
}

func UpdatePeerList(synchronizationChannelss elevatorConfig.SynchronizationChannels) {
	for {
		peerUpdate := <-synchronizationChannelss.PeerUpdateChannel

		alivePeersList = peerUpdate.Peers

		deadPeersList = append(deadPeersList, peerUpdate.Lost...)

		newPeer := peerUpdate.New

		if newPeer != "" {
			newList := []string{}

			for _, dead := range deadPeersList {
				if dead != newPeer {
					newList = append(newList, dead)
				}
			}

			deadPeersList = newList
		}
		fmt.Printf("Peers: %q\n", GetAlivePeersList())
		fmt.Printf("Dead: %q\n", GetDeadPeersList())
	}
}

func GetAlivePeersList() []string {
	return alivePeersList
}

func GetDeadPeersList() []string {
	return deadPeersList
}

// remove
func BroadcastElevatorWorldView(id string, BcastOutgoingMessagesChannel chan elevatorConfig.PeerView, elevatorSystem elevatorConfig.PeerView, elevatorChannel chan elevatorConfig.Elevator) {
	messageTimer := time.NewTimer(1 * time.Second) //should update slightly more often, maybe 30hz?
	for {

		//elevator := <-elevatorChannel

		//parse elevator to elevatorsystem (this function might fit here, might not)'
		//changed my mind, might have the parser run as its own event when recieveing elevators on the channel
		//and send them on the message channel... or this might ruin the switchcase
		//or this function might have become kind of useless unless we want to space the
		//world view brodcasts out a little more

		<-messageTimer.C

		BcastOutgoingMessagesChannel <- elevatorSystem

		messageTimer.Reset(1 * time.Second)
	}
}

//quite honestly, this might be useless, might just run these in the bigger sync function
//and have handlers that parse between types and set the info out on the correct channel

//naming suggestions for handlers
//HandleIncommingBroadcast
//HandleOutgoing Broadcast
//HandleIncommingOrderUpdate (should ths be in the state machine or mighgt be ein order protocol)

// remove
func RecieveBroadcastfWorldViewfFromPeer(BcastIncomingMessagesChannel chan elevatorConfig.PeerView) {

	for {
		incomingMessage := <-BcastIncomingMessagesChannel

		senderID := incomingMessage.OwnId
		/*
			direction := incomingMessage.Dirn
			request := incomingMessage.Requests
			behaviour := incomingMessage.Behaviour
		*/

		fmt.Println("Sender ID:", senderID)
		// Not elevatorsystem?
		//debuggingHelpers.Elevator_print(incomingMessage)

	}
}

func UpdateElevatorSystemFromElevator(elevator elevatorConfig.Elevator, peerView *elevatorConfig.PeerView) {
	SetBehavior(peerView, elevatorConfig.Behavior(elevator.Behavior))
	SetDirection(peerView, elevator.Direction)
	//fmt.Printf("Updating elevator system from elevator struct. Elevator floor: %d\n", elevator.Floor)
	if elevator.Floor >= 0 && elevator.Floor < elevatorConfig.N_FLOORS {
		SetFloor(peerView, elevator.Floor)
	}
}

/*
TODO: right, the state machine only operates on true/false from the requestsMatrix
There is a need for a function that transelates between the two OR the requestsMod
might need to be rewritten entirely. This might be in the order protocol mod
*/

// remove
func GetElevatorCabRequests(elevator elevatorConfig.Elevator) {

}
