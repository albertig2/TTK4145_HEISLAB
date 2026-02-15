
package driver



func main() {
    elevator_hardware_init() // add local IP å+ port when testing

    for {

	elevator_hardware_set_motor_direction(MDIR_DOWN);

	for(elevator_hardware_get_floor_sensor_signal() != 0) {}

	elevator_hardware_set_motor_direction(MDIR_UP);

	for(elevator_hardware_get_floor_sensor_signal() != 3) {}
    }
}