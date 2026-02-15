package driver

import (
	"fmt"
	"net"
	"sync"
)

var conn net.Conn
var socketmtx sync.Mutex

// Hardware constants and struckts
var N_floors int = 4
var N_btns int = 3

type BUTTON_TYPE int

const (
	BTN_UP   BUTTON_TYPE = 0
	BTN_DOWN             = 1
	BTN_CAB              = 2
)

type MOTOR_DIRECTION int

const (
	MDIR_DOWN MOTOR_DIRECTION = -1 //can aslo use iota here and remove the ints from the other
	MDIR_STOP                 = 0
	MDIR_UP                   = 1
)

func  elevator_hardware_init(address string) {

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

func elevator_hardware_set_motor_direction(dir MOTOR_DIRECTION) {
	socketmtx.Lock()
	defer socketmtx.Unlock()
	data := []byte{1, byte(dir), 0, 0}
	conn.Write(data)
}

func elevator_hardware_set_button_lamp(btn BUTTON_TYPE, floor int, value int) {
	if (floor >= 0) && (floor < N_floors) && (btn >= 0) && (int(btn) < N_btns) {

		socketmtx.Lock()
		defer socketmtx.Unlock()

		data := []byte{2, byte(btn), byte(floor), byte(value)}
		conn.Write(data)
	}
}

func elevator_hardware_set_floor_indicator(floor int) {

	if (floor >= 0) && (floor < N_floors) {
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
