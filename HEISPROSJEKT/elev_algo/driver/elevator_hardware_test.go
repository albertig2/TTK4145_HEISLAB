package driver

func test() {
	elevator_hardware_init("10.100.23.14:15657") // add local IP + port when testing

	for {

		elevator_hardware_set_motor_direction(D_DOWN)

		for elevator_hardware_get_floor_sensor_signal() != 0 {
		}

		elevator_hardware_set_motor_direction(D_UP)

		for elevator_hardware_get_floor_sensor_signal() != 3 {
		}
	}
}
