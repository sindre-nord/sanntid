package main

import (
	"elevator"
	"fmt"
	"settings"
)

var hall_calls [elevator.num_floors]bool                    // hall_calls[i] = true if there is a call at floor i
var elevator_states = make(map[int]*elevator.ElevatorState) // elevator_states[ID] = state for elevator #ID

func determineOrders(hall_calls elevator.Calls, states map[int]*elevator.ElevatorState) map[int]*elevator.Calls {
	_ = states
	orders := make(map[int]*elevator.Calls)
	orders[0] = &elevator.Calls{false, false, false, false}

	for i := range hall_calls {
		// If there is only one elevator, it services all orders
		orders[0][i] = hall_calls[i] || states[0].cab_calls[i]
	}
	return orders
}

// To ensure discovery of its own state at initialization, the elevator
// will start going upwards until it reaches any floor and then stop and
// enter normal operation
func discoverState(e *elevator.ElevatorState, ch_floor_passes chan int) {
	elevator.SetMotorDirection(e, elevator.elevio.MD_Up)
	for {
		select {
		case <-ch_floor_passes:
			elevator.SetMotorDirection(e, elevator.elevio.MD_Stop)
			return
		}
	}
}

func main() {
	elevator_states[0] = elevator.MakeElevatorState()
	self := elevator_states[0]
	hall_calls[2] = true
	elevator_states[0].cab_calls = elevator.Calls{false, true, false, true}
	// fmt.Println(determineOrders(hall_calls, elevator_states)[0])

	elevator.elevio.Init("localhost:15657", settings.NumFloors)

	ch_button_presses := make(chan elevator.elevio.ButtonEvent)
	ch_floor_passes := make(chan int)
	ch_obstructions := make(chan bool)
	ch_stop_pressed := make(chan bool)

	go elevator.elevio.PollButtons(ch_button_presses)
	go elevator.elevio.PollFloorSensor(ch_floor_passes)
	go elevator.elevio.PollObstructionSwitch(ch_obstructions)
	go elevator.elevio.PollStopButton(ch_stop_pressed)

	discoverState(self, ch_floor_passes)

	for {
		select {
		case press := <-ch_button_presses:
			fmt.Printf("Button %v, Floor %v\n", press.Button, press.Floor)
			elevator.SetMotorDirection(self, elevator.elevio.MD_Up)
		case floor := <-ch_floor_passes:
			fmt.Printf("Arrived at floor %v\n", floor)
		case is_obstructed := <-ch_obstructions:
			fmt.Printf("Is obstructed: %v\n", is_obstructed)
		case is_pressed := <-ch_stop_pressed:
			fmt.Printf("Stop button pressed: %v\n", is_pressed)
		}
	}
	return
}
