package communication

import (
	"HEISPROSJEKT/elevatorConfig"
	"Network-go/network/peers"
	"fmt"
	"strconv"
)

type messageStrc struct {
	Id      string
	Message string
}
type networkChannels struct {
	PeerUpdateChl                chan peers.PeerUpdate
	PeerTxEnableCh               chan bool
	BcastIncomingMessagesChannel chan messageStrc
	BcastOutgoingMessagesChannel chan messageStrc
}

var (
	alivePeersList []string
	deadPeersList  []string
)

func InitNetworkChannels() networkChannels {
	channels := networkChannels{
		PeerUpdateChl:                make(chan peers.PeerUpdate),
		PeerTxEnableCh:               make(chan bool),
		BcastIncomingMessagesChannel: make(chan messageStrc),
		BcastOutgoingMessagesChannel: make(chan messageStrc),
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

func BroadcastElevatorWorldView(id string, BcastOutgoingMessagesChannel chan messageStrc, harwdwareInputChannel chan elevatorConfig.Behavior) {
	//messageTimer := time.NewTimer(1*time.Second)
	for {
		behavior := <-harwdwareInputChannel

		outgoingMessage := messageStrc{
			Id:      id,
			Message: strconv.Itoa(int(behavior)),
		}
		fmt.Printf("SenderID: %s tried to send: %s \n", outgoingMessage.Id, outgoingMessage.Message)

		//<- messageTimer.C
		BcastOutgoingMessagesChannel <- outgoingMessage
		//messageTimer.Reset(1*time.Second)

	}

}

func RecieveBroadcastfWorldViewfFromPeer(BcastIncomingMessagesChannel chan messageStrc) {
	for {
		IncomingMessage := <-BcastIncomingMessagesChannel
		senderID := IncomingMessage.Id
		message := IncomingMessage.Message
		fmt.Printf("SenderID: %s said: %s \n", senderID, message)

	}

}
