package driver

import (
	"fmt"
	"net"
	"sync"
)

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

	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Init failed")

		return
	}

	data := []byte{0, 0, 0, 0}

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
