package main

import ("fmt")

var N_FLOORS int = 4
var N_BUTTONS int = 3

type Dirn int

const (
    D_Down Dirn = -1
    D_Stop Dirn = 0
    D_Up   Dirn = 1
)

type Button int

const (
    B_HallUp   Button = 0
    B_HallDown Button = 1
    B_Cab      Button = 2
)

type ElevatorBehaviour int

const (
    EB_Idle    ElevatorBehaviour = 0
    EB_DoorOpen ElevatorBehaviour = 1
    EB_Moving   ElevatorBehaviour = 2
)

type Elevator struct {
    floor      int
    dirn       Dirn
    requests   [N_Floor][N_Buttons]int
    behaviour  ElevatorBehaviour
    config     struct{ doorOpenDuration_s float64 }
}


func elevator_behaviorToString(eb ElevatorBehaviour) string{
    switch eb {
    case EB_Idle:
        return "EB_Idle"
    case EB_DoorOpen:
        return "EB_DoorOpen"
    case EB_Moving:
        return "EB_Moving"
    default:
        return "EB_UNDEFINED"
    }
}

func elevator_dirnToString(d Dirn) string {
    switch d {
    case D_Up: 
        return "D_Up"
    case D_Down:
        return "D_Down"
    case D_Stop:
        return "D_Stop"
    default:
        return "D_UNDEFINED"
    }
}

func elevator_buttonToString(b Button) string {
    switch b {
    case B_HallUp: 
        return "B_HallUp"
    case B_HallDown:
        return "B_HallDown"
    case B_Cab:
        return "B_Cab"
    default:
        return "B_UNDEFINED"
    }
}

func elevator_print(es Elevator) {
    fmt.Println("  +--------------------+")
    fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.floor)
    fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevator_dirnToString(es.dirn))
    fmt.Printf("  |%-6s = %-12.12s|\n", "behav", elevator_behaviorToString(es.behaviour))
    fmt.Println("  +--------------------+")
    fmt.Println("  |  | up  | dn  | cab |")
    for f := N_FLOORS - 1; f >= 0; i--{
        fmt.Printf("  | %d", f)
        for btn := 0; btn < N_BUTTONS; btn++{
            if (f == N_FLOORS - 1 && btn == B_HallUp) || (f == 0 && btn == B_HallDown) {
                fmt.Println("|     ")
            } else {
                if es.requests[f][btn] {
                    fmt.Println("|  #  ")
                } else {
                    fmt.Println("|  -  ")
                }
            }
        }
        fmt.Println("|")
    }
    fmt.Println("  +--------------------+")
}

func elevator_uninitialized() {
    elevator_hardware_init()
    return Elevator
}

func elevator_uninitialized() Elevator{
    elevator_hardware_init();
    es := Elevator{floor: -1, dirn: D_Stop, behaviour: EB_Idle, config: {doorOpenDuration_s: 3.0}}
    return es
}


/// LATER 
int elevator_floorSensor(void){
    return elevator_hardware_get_floor_sensor_signal();
}
int elevator_requestButton(int f, Button b){
    return elevator_hardware_get_button_signal((elevator_hardware_button_type_t)(b), f);
}
int elevator_stopButton(void){
    return elevator_hardware_get_stop_signal();
}
int elevator_obstruction(void){
    return elevator_hardware_get_obstruction_signal();
}

void elevator_floorIndicator(int f){
    elevator_hardware_set_floor_indicator(f);
}
void elevator_requestButtonLight(int f, Button b, int v){
    elevator_hardware_set_button_lamp((elevator_hardware_button_type_t)(b), f, v);
}
void elevator_doorLight(int v){
    elevator_hardware_set_door_open_lamp(v);
}
void elevator_stopButtonLight(int v){
    elevator_hardware_set_stop_lamp(v);
}
void elevator_motorDirection(Dirn d){
    elevator_hardware_set_motor_direction((elevator_hardware_motor_direction_t)(d));
}








