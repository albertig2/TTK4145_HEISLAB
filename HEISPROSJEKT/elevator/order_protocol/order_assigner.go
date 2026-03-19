package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronization"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type hallRequestAssignerPeerState struct {
	Behavior    string                        `json:"behavior"`
	Floor       int                           `json:"floor"`
	Direction   string                        `json:"direction"`
	CabRequests [elevatorConfig.N_FLOORS]bool `json:"cabRequests"`
}

type hallRequestAssignerPeerView struct {
	HallRequests [elevatorConfig.N_FLOORS][2]bool         `json:"hallRequests"`
	States       map[string]*hallRequestAssignerPeerState `json:"states"`
}

func buildHallRequestAssignerPeerView(peerView elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition) hallRequestAssignerPeerView {
	hallRequestAssignerPeerView := hallRequestAssignerPeerView{
		HallRequests: [elevatorConfig.N_FLOORS][2]bool{},
		States:       make(map[string]*hallRequestAssignerPeerState),
	}

	for _, peerId := range peerView.AlivePeers {
		peerState := peerView.States[peerId]
		hallRequestAssignerPeerView.States[peerId] = &hallRequestAssignerPeerState{
			Behavior:    elevatorConfig.BehaviorToString(peerState.Behavior),
			Floor:       peerState.Floor,
			Direction:   elevatorConfig.DirectionToString(peerState.Direction),
			CabRequests: [elevatorConfig.N_FLOORS]bool{},
		}
	}

	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDirection := range synchronization.HallDirections {
			if hallRequestTransitions[floor][hallDirection] == PendingToAssigned {
				hallRequestAssignerPeerView.HallRequests[floor][hallDirection] = true
			}
		}
	}
	return hallRequestAssignerPeerView
}

func hallRequestAssigner(peerView *elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition) map[string][][2]bool {
	Executable := ""
	switch runtime.GOOS {
	case "linux":
		Executable = "hall_request_assigner"
	case "windows":
		Executable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := buildHallRequestAssignerPeerView(*peerView, hallRequestTransitions)

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
	output := new(map[string][][2]bool)
	error = json.Unmarshal(outputBytes, &output)
	if error != nil {
		fmt.Println("json.Unmarshal error: ", error)
		return nil
	}
	return *output
}
