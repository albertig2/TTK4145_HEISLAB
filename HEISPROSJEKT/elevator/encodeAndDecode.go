package main

import "encoding/json"

func encodeJson(system *ElevatorSystem) string {
	jsonData, _ := json.Marshal(&system)
	return string(jsonData)
}

func decodeJson(jsonStr string) ElevatorSystem {
	var system ElevatorSystem
	json.Unmarshal([]byte(jsonStr), &system)
	return system
}
