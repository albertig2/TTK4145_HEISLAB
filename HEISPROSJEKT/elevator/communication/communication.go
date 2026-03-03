package communication

import (
	"Network-go/network/peers"
	"strconv"	
)

type networkChannels struct {
	PeerUpdateChl chan peers.PeerUpdate
	PeerTxEnableCh chan bool
}

var (
	alivePeersList []string
	deadPeersList []string
)

func InitNetworkChannels() networkChannels{
	channels := networkChannels{
		PeerUpdateChl: make(chan peers.PeerUpdate),
		PeerTxEnableCh: make(chan bool),
	}

	return channels
}

func StartPeerNetworking(port int, id int, channels networkChannels) {
	go peers.Receiver(port, channels.PeerUpdateChl)
	go peers.Transmitter(port, strconv.Itoa(id), channels.PeerTxEnableCh)
}



func UpdatePeerList(channels networkChannels){
	peerUpdate := <- channels.PeerUpdateChl

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

}

func GetAlivePeersList() []string {
	return alivePeersList
}

func GetDeadPeersList() []string {
	return deadPeersList
}



