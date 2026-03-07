package elevatorConfig

const N_FLOORS int = 4
const N_BUTTONS int = 3

type Direction int

const (
	Down Direction = -1
	Stop Direction = 0
	Up   Direction = 1
)

type Button int

const (
	HallUp   Button = 0
	HallDown Button = 1
	Cab      Button = 2
)

func DirectionToString(direction Direction) string {
	switch direction {
	case Up:
		return "up"
	case Down:
		return "down"
	case Stop:
		return "stop"
	default:
		return "undefined"
	}
}

func ButtonToString(b Button) string {
	switch b {
	case HallUp:
		return "HallUp"
	case HallDown:
		return "HallDown"
	case Cab:
		return "Cab"
	default:
		return "Undefined"
	}
}

/*


type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2
)

type Config struct {
	doorOpenDuration_s time.Duration
}

type Elevator struct {
	floor     int
	direction	requests  [N_FLOORS][N_BUTTONS]bool
	behaviour ElevatorBehaviour
	config    Config
}
*/
/*
func Elevator_behaviorToString(eb ElevatorBehaviour) string {
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




*/
