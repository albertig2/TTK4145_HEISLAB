package debuggingHelpers

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
)

func Elevator_print(es elevatorConfig.Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.Floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevatorConfig.DirectionToString(es.Direction))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", elevatorConfig.BehaviorToString(es.Behavior))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevatorConfig.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if (f == elevatorConfig.N_FLOORS-1 && elevatorConfig.Button(btn) == elevatorConfig.HallUp) || (f == 0 && elevatorConfig.Button(btn) == elevatorConfig.HallDown) {
				fmt.Printf("|     ")
			} else {
				if es.Requests[f][btn] {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}
