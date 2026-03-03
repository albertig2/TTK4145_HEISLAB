package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

func hallRequestAssigner(system *ElevatorSystem) {
	Executable := ""
	switch runtime.GOOS {
	case "linux":
		Executable = "hall_request_assigner"
	case "windows":
		Executable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := system

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}
	ret, err := exec.Command("../cost_fns/hall_request_assigner/"+Executable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for id, hallRequests := range *output {
		fmt.Printf("%6v :  %+v\n", id, hallRequests)
	}

	for id, hallRequests := range *output {
		for floor, hallRequest := range hallRequests {
			if hallRequest[0] || hallRequest[1] {
				intId, _ := strconv.Atoi(id)
				setCabRequests(system, intId, floor, true)
			}
		}
	}
}
