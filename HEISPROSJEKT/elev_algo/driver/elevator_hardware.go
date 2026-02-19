package driver

import (
	"fmt"
	"net"
	"sync"
)

var N_FLOORS int = 4
var N_BUTTONS int = 3

var conn net.Conn
var socketmtx sync.Mutex

// Hardware constants and struckts
var number_floors = 4

type elevator_hardware_motor_direction_t int
type elevator_hardware_button_type_t int

const (
	B_HallUP   elevator_hardware_button_type_t = 0
	B_HallDown                                 = 1
	B_Cab                                      = 2
)

const (
	D_DOWN elevator_hardware_motor_direction_t = -1 //can aslo use iota here and remove the ints from the other
	D_STOP                                     = 0
	D_UP                                       = 1
)

func elevator_hardware_init(address string) {

	var err error = nil

	conn, err = net.Dial("tcp", address)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Init failed")

		return
	}

	data := []byte{0, 0, 0, 0}

	conn.Write(data)

}

func elevator_hardware_set_motor_direction(dir elevator_hardware_motor_direction_t) {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	data := []byte{1, byte(dir), 0, 0}
	conn.Write(data)
}

func elevator_hardware_set_button_lamp(btn elevator_hardware_button_type_t, floor int, value int) {
	if (floor >= 0) && (floor < N_FLOORS) && (btn >= 0) && (int(btn) < N_BUTTONS) {

		socketmtx.Lock()
		defer socketmtx.Unlock()

		data := []byte{2, byte(btn), byte(floor), byte(value)}
		conn.Write(data)
	}
}

func elevator_hardware_set_floor_indicator(floor int) {

	if (floor >= 0) && (floor < N_FLOORS) {
		socketmtx.Lock()
		defer socketmtx.Unlock()
		data := []byte{3, byte(floor), 0, 0}
		conn.Write(data)
	}

}

func elevator_hardware_set_door_open_lamp(value int) {

	socketmtx.Lock()
	defer socketmtx.Unlock()
	data := []byte{4, byte(value), 0, 0}
	conn.Write(data)
}

func elevator_hardware_set_stop_lamp(value int) {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	conn.Write([]byte{5, byte(value), 0, 0})
}

func elevator_hardware_get_button_signal(button elevator_hardware_button_type_t, floor int) bool {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	conn.Write([]byte{6, byte(button), byte(floor), 0})
	buf := make([]byte, 4)
	conn.Read(buf)
	return buf[1] != 0
}

func elevator_hardware_get_floor_sensor_signal() int {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	conn.Write([]byte{7, 0, 0, 0})
	buf := make([]byte, 4)
	conn.Read(buf)
	if buf[1] != 0 {
		return int(buf[2])
	} else {
		return -1
	}
}

func elevator_hardware_get_stop_signal() bool {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	conn.Write([]byte{8, 0, 0, 0})
	buf := make([]byte, 4)
	conn.Read(buf)
	return buf[1] != 0
}

func elevator_hardware_get_obstruction_signal() bool {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	conn.Write([]byte{9, 0, 0, 0})
	buf := make([]byte, 4)
	conn.Read(buf)
	return buf[1] != 0
}
