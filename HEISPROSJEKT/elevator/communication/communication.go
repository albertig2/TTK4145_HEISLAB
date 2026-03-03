package communication

import (
	"Network-go/network/peers"
	"strconv"	
)
//privte alivelist
//private deadlist
//possebly somthing with new elevators

//struct for communication channels
type networkChannels struct {
	PeerUpdateChl chan peers.PeerUpdate
	PeerTxEnableCh chan bool
}

var (
	alivePeersList []string
	deadPeersList []string
)

func InitNetworkChannels(id int, port int ) networkChannels{
	channels := networkChannels{
		PeerUpdateChl: make(chan peers.PeerUpdate),
		PeerTxEnableCh: make(chan bool),
	}

	go peers.Receiver(port, channels.PeerUpdateChl)
	go peers.Transmitter(port, strconv.Itoa(id), channels.PeerTxEnableCh)

	return channels
}



// func maintinElevatorNetworkStatus(peerTxenable, peerUpdateChl)
	//peerUpdate := <- peerUpdateChl
	//alive list = peerupdate.pers 
	//deadList = peerUpdate.deadlist

func UpdatePeerList(channels networkChannels){
	peerUpdate := <- channels.PeerUpdateChl
	alivePeersList = peerUpdate.Peers
	deadPeersList = peerUpdate.Lost

}

func GetAlivePeersList() []string {
	return alivePeersList
}

func GetDeadPeersList() []string {
	return deadPeersList
}



