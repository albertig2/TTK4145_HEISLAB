package communication

import (
	"Network-go/network/peers"	
)
//privte alivelist
//private deadlist
//possebly somthing with new elevators

//struct for communication channels
type networkChannels struct {

}

func initNetworkChannels(id int, port int ){

}

// func maintinElevatorNetworkStatus(peerTxenable, peerUpdateChl)
	//peerUpdate := <- peerUpdateChl
	//alive list = peerupdate.pers 
	//deadList = peerUpdate.deadlist

func getAlivePeersList() []string {
	return aliveList
}

func getDeadPeersList() []string {
	return deadList
}



