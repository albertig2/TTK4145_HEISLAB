package communication

import (
	//"HEISPROSJEKT/hardware"
	"HEISPROSJEKT/elevatorHardware"
	"Network-go/network/peers"
	"fmt"
	"strconv"
	"time"
)

type messageStrc struct {
	Id      string
	Message string
}
type networkChannels struct {
	PeerUpdateChl                chan peers.PeerUpdate
	PeerTxEnableCh               chan bool
	BcastIncomingMessagesChannel chan elevatorHardware.Elevator
	BcastOutgoingMessagesChannel chan elevatorHardware.Elevator
}

var (
	alivePeersList []string
	deadPeersList  []string
)

func InitNetworkChannels() networkChannels {
	channels := networkChannels{
		PeerUpdateChl:                make(chan peers.PeerUpdate),
		PeerTxEnableCh:               make(chan bool),
		BcastIncomingMessagesChannel: make(chan elevatorHardware.Elevator),
		BcastOutgoingMessagesChannel: make(chan elevatorHardware.Elevator),
	}

	return channels
}

func StartPeerNetworking(port int, id int, channels networkChannels) {
	fmt.Println("start")

	go peers.Receiver(port, channels.PeerUpdateChl)
	go peers.Transmitter(port, strconv.Itoa(id), channels.PeerTxEnableCh)
}

func UpdatePeerList(channels networkChannels) {
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

func BroadcastElevatorWorldView(id string, BcastOutgoingMessagesChannel chan elevatorHardware.Elevator, elevatorChannel chan elevatorHardware.Elevator) {
	messageTimer := time.NewTimer(1 * time.Second)
	for {

		elevator := <- elevatorChannel

		<-messageTimer.C
		
		BcastOutgoingMessagesChannel <- elevator

		messageTimer.Reset(1 * time.Second)

	}
}

func RecieveBroadcastfWorldViewfFromPeer(BcastIncomingMessagesChannel chan elevatorHardware.Elevator) {
	for {
		incomingMessage := <-BcastIncomingMessagesChannel

		senderID := incomingMessage.OwnId
		/*
			direction := incomingMessage.Dirn
			request := incomingMessage.Requests
			behaviour := incomingMessage.Behaviour
		*/

		fmt.Println("Sender ID:", senderID)
		elevatorHardware.Elevator_print(incomingMessage)

	}

}
