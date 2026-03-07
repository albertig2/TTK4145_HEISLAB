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

type Behavior int

const (
	Idle     Behavior = 0
	DoorOpen Behavior = 1
	Moving   Behavior = 2
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
		return "hallUp"
	case HallDown:
		return "hallDown"
	case Cab:
		return "cab"
	default:
		return "undefined"
	}
}

func BehaviorToString(eb Behavior) string {
	switch eb {
	case Idle:
		return "idle"
	case DoorOpen:
		return "doorOpen"
	case Moving:
		return "moving"
	default:
		return "undefined"
	}
}

/*

type Config struct {
	doorOpenDuration_s time.Duration
}

type Elevator struct {
	floor     int
	direction	requests  [N_FLOORS][N_BUTTONS]bool
	behavior  ElevatorBehavior
	config    Config
}
*/
/*





 */
