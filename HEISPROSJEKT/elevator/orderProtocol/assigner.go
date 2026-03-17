package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	//"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/synchronization"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type BoolElevatorState struct {
	Behavior    string                        `json:"behavior"`
	Floor       int                           `json:"floor"`
	Direction   string                        `json:"direction"`
	CabRequests [elevatorConfig.N_FLOORS]bool `json:"cabRequests"`
}

type BoolElevatorSystem struct {
	HallRequests [elevatorConfig.N_FLOORS][2]bool `json:"hallRequests"`
	States       map[string]*BoolElevatorState    `json:"states"`
}

// Converts ElevatorSystem and order status to a boolean-based system for assignment logic
func BuildBoolElevatorSystem(system elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition) BoolElevatorSystem {
	boolSystem := BoolElevatorSystem{
		HallRequests: [elevatorConfig.N_FLOORS][2]bool{},
		States:       make(map[string]*BoolElevatorState),
	}

	for _, peerId := range system.AlivePeers {
		idState := system.States[peerId]

		boolSystem.States[peerId] = &BoolElevatorState{
			Behavior:    elevatorConfig.BehaviorToString(idState.Behavior),
			Floor:       idState.Floor,
			Direction:   elevatorConfig.DirectionToString(idState.Direction),
			CabRequests: [elevatorConfig.N_FLOORS]bool{},
		}
	}

	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDir := range synchronization.HallDirections {
			if hallRequestTransitions[floor][hallDir] == PendingToAssigned {
				boolSystem.HallRequests[floor][hallDir] = true
			}
		}
	}
	return boolSystem
}

func HallRequestAssigner(system *elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition) map[string][][2]bool {
	Executable := ""
	switch runtime.GOOS {
	case "linux":
		Executable = "hall_request_assigner"
	case "windows":
		Executable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := BuildBoolElevatorSystem(*system, hallRequestTransitions)

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return nil
	}
	ret, err := exec.Command("../cost_fns/hall_request_assigner/"+Executable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return nil
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return nil
	}
	/*
		fmt.Printf("output: \n")
		for id, hallRequests := range *output {
			fmt.Printf("%6v :  %+v\n", id, hallRequests)
		}
	*/
	return *output
}
