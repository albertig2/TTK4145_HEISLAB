package synchronisation

import (
	"HEISPROSJEKT/debuggingHelpers"
	"HEISPROSJEKT/elevatorConfig"
	"Network-go/network/peers"
	
	"fmt"
	"strconv"
	"time"
)


type synchronisationChannels struct {
	PeerUpdateChl                chan peers.PeerUpdate
	PeerTxEnableCh               chan bool
	BcastIncomingMessagesChannel chan ElevatorSystem
	BcastOutgoingMessagesChannel chan ElevatorSystem
	UpdateElevatorSystem chan ElevatorSystem
	//new order channel?
}

var (
	alivePeersList []string
	deadPeersList  []string
)

func InitNetworkChannels () synchronisationChannels {
	channels := synchronisationChannels{
		PeerUpdateChl:                make(chan peers.PeerUpdate),
		PeerTxEnableCh:               make(chan bool),
		BcastIncomingMessagesChannel: make(chan ElevatorSystem),
		BcastOutgoingMessagesChannel: make(chan ElevatorSystem),
		UpdateElevatorSystem: make(chan ElevatorSystem),
	}

	return channels
}

func StartPeerNetworking(port int, id int, channels synchronisationChannels) {
	fmt.Println("start")

	go peers.Receiver(port, channels.PeerUpdateChl)
	go peers.Transmitter(port, strconv.Itoa(id), channels.PeerTxEnableCh)
}

func UpdatePeerList(channels synchronisationChannels) {
	for {
		peerUpdate := <-channels.PeerUpdateChl

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



func BroadcastElevatorWorldView(id string, BcastOutgoingMessagesChannel chan ElevatorSystem, elevatorSystem ElevatorSystem, elevatorChannel chan elevatorConfig.Elevator) {
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

func RecieveBroadcastfWorldViewfFromPeer(BcastIncomingMessagesChannel chan elevatorConfig.Elevator) {

	for {
		incomingMessage := <-BcastIncomingMessagesChannel

		senderID := incomingMessage.OwnId
		/*
			direction := incomingMessage.Dirn
			request := incomingMessage.Requests
			behaviour := incomingMessage.Behaviour
		*/

		fmt.Println("Sender ID:", senderID)
		debuggingHelpers.Elevator_print(incomingMessage)

	}
}

func UpdateElevatorSystemFromELevator(elevator elevatorConfig.Elevator, elevatorSystem *ElevatorSystem) {

	SetBehavior(elevatorSystem, elevatorConfig.Behavior(elevator.Direction))
	SetDirection(elevatorSystem, elevator.Direction)
	SetFloor(elevatorSystem, elevator.Floor)

}
/*
TODO: right, the state machine only operates on true/false from the requestsMatrix
There is a need for a function that transelates between the two OR the requestsMod
might need to be rewritten entirely. This might be in the order protocol mod 

*/
func GetElevatorCabRequests(elevator elevatorConfig.Elevator){

}
