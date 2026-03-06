package communication

import "encoding/json"

func encodeElevatorSystem(system *ElevatorSystem) string {
	jsonData, _ := json.Marshal(&system)
	return string(jsonData)
}

func decodeElevatorSystem(jsonStr string) ElevatorSystem {
	var system ElevatorSystem
	json.Unmarshal([]byte(jsonStr), &system)
	return system
}
