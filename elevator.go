package elevator

import (
	"elevator_controller/src/elevio"
	"settings"
)

type Calls = [settings.NumFloors]bool

// ElevatorState is a struct that contains the state of an elevator
// It is not ment to be initialized directly, but rather through the
// Make function MakeElevatorState
type ElevatorState struct {
	last_passed_floor uint		// Last floor passed by the elevator
	direction int 				// -1 = down, 0 = stationary, 1 = up
	cab_calls [settings.NumFloors]bool 	// cab_calls[i] = true if elevator should stop at floor i
}

// MakeElevatorState is a constructor for the ElevatorState struct.
// It initializes the struct with default values, and theese are the only 
// values that should be used when initializing a new ElevatorState
func MakeElevatorState() *ElevatorState{
	e := ElevatorState{0, elevio.MD_Stop, [settings.NumFloors]bool{false, false, false, false}}
	return &e
}

func SetMotorDirection(e *ElevatorState, dir int) {
	e.direction = dir
	elevio.SetMotorDirection(dir)
}
