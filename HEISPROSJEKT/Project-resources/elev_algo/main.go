package main

import ("fmt")


func main () {
    fmt.Println("Started!")
    var elevator Elevator = elevator_uninitialized()
    var inputPollRate_ms int = 25

    con_load("elevator.con",
        con_val("doorOpenDuration_s", &elevator.config.doorOpenDuration_s, "%lf")
        con_val("inputPollRate_ms", &inputPollRate_ms, "%d")
    )
    
    if elevator_floorSensor() == -1 {
        fsm_onInitBetweenFloors(&elevator);
    }
    
    var prevButtons [N_FLOORS][N_BUTTONS]bool
    var prevFloor int = -1
    for true {
        { // Request button
            for f := 0; f < N_FLOORS; f++ {
                for b := 0; b < N_BUTTONS; b++ {
                    var v bool = elevator_requestButton(f, b)
                    if v  &&  v != prevButtons[f][b] {
                        fsm_onRequestButtonPress(&elevator, f, b)
                    }
                    prevButtons[f][b] = v
                }
            }
        }
        
        { // Floor sensor
            var f int = elevator_floorSensor();
            if f != -1  &&  f != prevFloor {
                fsm_onFloorArrival(&elevator, f);
            }
            prevFloor = f
        }
        
        
        { // Timer
            if timer_timedOut() {
                timer_stop()
                fsm_onDoorTimeout(&elevator)
            }
        }
        
        usleep(inputPollRate_ms*1000)
    }

}