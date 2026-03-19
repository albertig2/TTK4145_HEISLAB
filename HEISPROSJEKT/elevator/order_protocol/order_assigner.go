package orderProtocol

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"HEISPROSJEKT/synchronization"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type hallOrderAssignerPeerState struct {
	Behavior  string                              `json:"behavior"`
	Floor     int                                 `json:"floor"`
	Direction string                              `json:"direction"`
	CabOrders [elevatorConfig.NumberOfFloors]bool `json:"cabRequests"`
}

type hallOrderAssignerPeerView struct {
	HallOrders [elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]bool `json:"hallRequests"`
	States     map[string]*hallOrderAssignerPeerState                                  `json:"states"`
}

func buildHallOrderAssignerPeerView(peerView elevatorConfig.PeerView, hallOrderTransitions [elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]OrderTransition) hallOrderAssignerPeerView {
	hallOrderAssignerPeerView := hallOrderAssignerPeerView{
		HallOrders: [elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]bool{},
		States:     make(map[string]*hallOrderAssignerPeerState),
	}

	for _, peerId := range peerView.AlivePeers {
		peerState := peerView.States[peerId]
		hallOrderAssignerPeerView.States[peerId] = &hallOrderAssignerPeerState{
			Behavior:  elevatorConfig.BehaviorToString(peerState.Behavior),
			Floor:     peerState.Floor,
			Direction: elevatorConfig.DirectionToString(peerState.Direction),
			CabOrders: [elevatorConfig.NumberOfFloors]bool{},
		}
	}

	for floor := range elevatorConfig.NumberOfFloors {
		for _, hallDirection := range synchronization.HallDirections {
			if hallOrderTransitions[floor][hallDirection] == PendingToAssigned {
				hallOrderAssignerPeerView.HallOrders[floor][hallDirection] = true
			}
		}
	}
	return hallOrderAssignerPeerView
}

func hallOrderAssigner(peerView *elevatorConfig.PeerView, hallOrderTransitions [elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]OrderTransition) map[string][][elevatorConfig.NumberOfHallButtons]bool {
	Executable := ""
	switch runtime.GOOS {
	case "linux":
		Executable = "hall_request_assigner"
	case "windows":
		Executable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := buildHallOrderAssignerPeerView(*peerView, hallOrderTransitions)

	jsonBytes, error := json.Marshal(input)
	if error != nil {
		fmt.Println("json.Marshal error: ", error)
		return nil
	}
	outputBytes, error := exec.Command("../cost_fns/hall_request_assigner/"+Executable, "-i", string(jsonBytes)).CombinedOutput()
	if error != nil {
		fmt.Println("exec.Command error: ", error)
		fmt.Println(string(outputBytes))
		return nil
	}
	output := new(map[string][][elevatorConfig.NumberOfHallButtons]bool)
	error = json.Unmarshal(outputBytes, &output)
	if error != nil {
		fmt.Println("json.Unmarshal error: ", error)
		return nil
	}
	return *output
}
