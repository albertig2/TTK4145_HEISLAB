package synchronization

import (
	"HEISPROSJEKT/elevatorConfig"
	"encoding/json"
	//"HEISPROSJEKT/elevatorHardware"
)

func EncodeElevatorSystem(peerView *elevatorConfig.PeerView) string {
	jsonData, _ := json.Marshal(&peerView)
	return string(jsonData)
}

func DecodeElevatorSystem(jsonStr string) elevatorConfig.PeerView {
	var peerView elevatorConfig.PeerView
	json.Unmarshal([]byte(jsonStr), &peerView)
	return peerView
}
