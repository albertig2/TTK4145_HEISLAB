package communication

import( 
	"encoding/json"
	"HEISPROSJEKT/elevatorHardware"
	

)

func EncodeElevatorSystem(system *elevatorHardware.ElevatorSystem) string {
	jsonData, _ := json.Marshal(&system)
	return string(jsonData)
}

func DecodeElevatorSystem(jsonStr string) elevatorHardware.ElevatorSystem {
	var system elevatorHardware.ElevatorSystem
	json.Unmarshal([]byte(jsonStr), &system)
	return system
}
