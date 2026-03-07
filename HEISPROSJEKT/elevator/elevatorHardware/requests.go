package elevatorHardware

type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}

// static i C -> privat i Go (liten forbokstav)
func Requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

// “API” lik headeren: behold navnet
func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case D_Up:
		if Requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		}
		if Requests_here(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		}
		if Requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
		return DirnBehaviourPair{D_Stop, EB_Idle}

	case D_Down:
		if Requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
		if Requests_here(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		}
		if Requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		}
		return DirnBehaviourPair{D_Stop, EB_Idle}

	case D_Stop: // samme kommentar som i C
		if Requests_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		}
		if Requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		}
		if Requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
		return DirnBehaviourPair{D_Stop, EB_Idle}

	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}

func Requests_shouldStop(e Elevator) bool {
	switch e.Dirn {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] ||
			e.Requests[e.Floor][B_Cab] ||
			!Requests_below(e)

	case D_Up:
		return e.Requests[e.Floor][B_HallUp] ||
			e.Requests[e.Floor][B_Cab] ||
			!Requests_above(e)

	case D_Stop:
		fallthrough
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) bool {
	return e.Floor == btn_floor &&
		((e.Dirn == D_Up && btn_type == B_HallUp) ||
			(e.Dirn == D_Down && btn_type == B_HallDown) ||
			e.Dirn == D_Stop ||
			btn_type == B_Cab)
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.Requests[e.Floor][B_Cab] = false

	switch e.Dirn {
	case D_Up:
		if !Requests_above(e) && !e.Requests[e.Floor][B_HallUp] {
			e.Requests[e.Floor][B_HallDown] = false
		}
		e.Requests[e.Floor][B_HallUp] = false

	case D_Down:
		if !Requests_below(e) && !e.Requests[e.Floor][B_HallDown] {
			e.Requests[e.Floor][B_HallUp] = false
		}
		e.Requests[e.Floor][B_HallDown] = false

	case D_Stop:
		fallthrough
	default:
		e.Requests[e.Floor][B_HallUp] = false
		e.Requests[e.Floor][B_HallDown] = false
	}

	return e
}
