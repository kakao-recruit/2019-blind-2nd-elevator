package elevator

import (
	"fmt"
	"2019-blind-2nd-elevator/elevator/internal/mydb"
)

type Elevator struct {
	Id int `json:"id"`
	Floor int `json:"floor"`
	Passengers mydb.Calls `json:"passengers"`
	Status string `json:"status"`

	MaxFloor int `json:"-"`
	MaxPassengers int `json:"-"`
}

func NewElevator(n, maxFloor, maxPassengers int) Elevator {
	return Elevator {
		n, 1, make(mydb.Calls, 0, maxPassengers), "STOPPED", maxFloor, maxPassengers,
	}
}

func (e *Elevator) Up() error {
	if e.Status != "STOPPED" && e.Status != "UPWARD" {
		return fmt.Errorf("Wrong action")
	}

	if e.Floor < e.MaxFloor {
		e.Floor += 1
	}

	e.Status = "UPWARD"
	return nil
}

func (e *Elevator) Down() error {
	if e.Status != "STOPPED" && e.Status != "DOWNWARD" {
		return fmt.Errorf("Wrong action")
	}

	if e.Floor > 1 {
		e.Floor -= 1
	}

	e.Status = "DOWNWARD"
	return nil
}

func (e *Elevator) Stop() error {
	if e.Status != "UPWARD" && e.Status != "DOWNWARD" && e.Status != "STOPPED" {
		return fmt.Errorf("Wrong action")
	}

	e.Status = "STOPPED"
	return nil
}

func (e *Elevator) Open() error {
	if e.Status != "STOPPED" && e.Status != "OPENED" {
		return fmt.Errorf("Wrong action")
	}

	e.Status = "OPENED"
	return nil
}

func (e *Elevator) Close() error {
	if e.Status != "OPENED" {
		return fmt.Errorf("Wrong action")
	}

	e.Status = "STOPPED"
	return nil
}

func (e *Elevator) IsFull() bool {
	return len(e.Passengers) >= e.MaxPassengers
}

func (e *Elevator) Enter(call mydb.Call) error {
	if e.Status != "OPENED" {
		return fmt.Errorf("Wrong action")
	}

	if e.IsFull() {
		return fmt.Errorf("Exceeding the capacity of the elevator")
	}

	if call.Start != e.Floor {
		return fmt.Errorf("No passenger for %d at floor %d", call.Id, e.Floor)
	}

	e.Passengers = append(e.Passengers, call)
	return nil
}

func (e *Elevator) Exit(id int) (*mydb.Call, error) {
	if e.Status != "OPENED" {
		return nil, fmt.Errorf("Wrong action")
	}

	var found int = -1
	for i, v := range e.Passengers {
		if v.Id == id {
			found = i
			break
		}
	}

	if found < 0 {
		return nil, fmt.Errorf("No passenger for %d in the elevator", id)
	}

	p := e.Passengers[found]
	e.Passengers = append(e.Passengers[:found], e.Passengers[found+1:]...)

	return &p, nil
}
