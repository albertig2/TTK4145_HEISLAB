package driver

import (
	 "os"
	"math/rand"
	"net"
	"fmt"
	"sync"

)


var sockfd int
var socketmtx sync.Mutex

//Hardware constants and struckts
var number_floors = 4;

type BUTTON_TYPE int
type MOTOR_DIRECTION int


const (
	BTN_UP BUTTON_TYPE = 0
	BTN_DOWN = 1
	BTN_CAB = 2
)

const (
	MDIR_DOWN MOTOR_DIRECTION= -1 //can aslo use iota here and remove the ints from the other 
	MDIR_STOP = 0
	MDIR_UP  = 1
	
)


func elevator_hardware_init( address string){

	
	conn, err := net.Dial("tcp", address)

	if err != nil{
		fmt.Println(err)
		fmt.Println("Init failed")
		
		return
	}


	data := []byte{0,0,0,0}

	conn.Write(data)

}
