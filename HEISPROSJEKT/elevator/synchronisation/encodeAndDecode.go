package synchronisation

import (
	"HEISPROSJEKT/elevatorConfig"
	"encoding/json"
	//"HEISPROSJEKT/elevatorHardware"
)

func EncodeElevatorSystem(system *elevatorConfig.ElevatorSystem) string {
	jsonData, _ := json.Marshal(&system)
	return string(jsonData)
}

func DecodeElevatorSystem(jsonStr string) elevatorConfig.ElevatorSystem {
	var system elevatorConfig.ElevatorSystem
	json.Unmarshal([]byte(jsonStr), &system)
	return system
}
