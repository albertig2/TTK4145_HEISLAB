package synchronisation

import( 
	"encoding/json"
	//"HEISPROSJEKT/elevatorHardware"

	

)

func EncodeElevatorSystem(system *ElevatorSystem) string {
	jsonData, _ := json.Marshal(&system)
	return string(jsonData)
}

func DecodeElevatorSystem(jsonStr string) ElevatorSystem {
	var system ElevatorSystem
	json.Unmarshal([]byte(jsonStr), &system)
	return system
}
