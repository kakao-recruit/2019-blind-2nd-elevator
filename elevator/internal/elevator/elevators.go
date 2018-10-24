package elevator

import (
	"os"
	"path"
	"math/rand"
	"time"
	"fmt"
	"2019-blind-2nd-elevator/elevator/internal/mydb"
	"github.com/spf13/viper"
	eval "2019-blind-2nd-elevator/elevator/internal/eval"
	log "github.com/sirupsen/logrus"
)

type Elevators struct {
	Token string `json:"token"`
	Timestamp int `json:"timestamp"`

	Cars []Elevator `json:"cars"`

	Problem int `json:"-"`
	Calls mydb.Calls `json:"-"`
	LastCallTs int `json:"-"`

	Inputs mydb.Calls `json:"-"`
	Outputs mydb.Outputs `json:"-"`
	Logs []string `json:"-"`

	done bool 
}

var (
	length = 5
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateToken() string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewElevators(problem, carCount int) (*Elevators, error) {
	token := generateToken()

	maxPassengers := viper.GetInt("MaxPassengers")

	calls := mydb.GetCalls(problem)

	maxFloor := 0

	switch problem {
	case 0:
		maxFloor = 5
	case 1:
		maxFloor = 25
	case 2:
		maxFloor = 25
	}

	cars := make([]Elevator, 0)
	for i := 0; i < carCount; i++ {
		cars = append(cars, NewElevator(i, maxFloor, maxPassengers))
	}

	inputs := make(mydb.Calls, len(calls))
	copy(inputs, calls)

	return &Elevators{token, 0, cars, problem, calls, -1, inputs, nil, nil, false}, nil
}

func (e *Elevators) Logging(action string, car int) {
	status, passenger := "", ""

	switch e.Cars[car].Status {
	case "UPWARD":
		status = "U"
	case "DOWNWARD":
		status = "D"
	case "OPENED":
		status = "O"
	case "STOPPED":
		status = "S"
	}

	for _, p := range e.Cars[car].Passengers {
		if passenger == "" {
			passenger = fmt.Sprintf("%d", p.End)
		} else {
			passenger = fmt.Sprintf("%s:%d", passenger, p.End)
		}
	}

	log := fmt.Sprintf("e,%d,%d,%d,%s,%s,%s", e.Timestamp, car, e.Cars[car].Floor, action, status, passenger)

	e.Logs = append(e.Logs, log)
}

func (e *Elevators) Tick() {
	e.Timestamp += 1
}

func (e *Elevators) Count() int {
	return len(e.Cars)
}

func (e *Elevators) IsEnd() bool {
	if len(e.Calls) != 0 {
		return false
	}

	for i := range e.Cars {
		if len(e.Cars[i].Passengers) != 0 {
			return false
		}
	}

	return true
}

func (e *Elevators) Done() {
	if e.done {
		return
	}

	e.done = true

	wait, travel, total, lastTs, err := e.Evaluate()
	var status string
	if err != nil {
		status = err.Error()
	} else {
		status = "OK"
	}

	dir := viper.GetString("LogDir")
	path := path.Join(dir, fmt.Sprintf("%s.log", e.Token))

	f, err := os.Create(path)
	if err != nil {
		log.Debug("Failed to create log file")
		return
	}

	defer f.Close()

	f.WriteString(fmt.Sprintf("%s %s\n", e.Token, status))
	if status == "OK" {
		f.WriteString(fmt.Sprintf("%f %f %f %d\n", wait, travel, total, lastTs))
	}

	for _, l := range e.Logs {
		f.WriteString(l + "\n")
	}

	f.Sync()
}

func (e *Elevators) Evaluate() (float64, float64, float64, int, error) {
	return eval.Evaluate(e.Inputs, e.Outputs)
}

func (e *Elevators) Up(car int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if err := e.Cars[car].Up(); err != nil {
		return err
	}

	e.Logging("U", car)
	return nil
}

func (e *Elevators) Down(car int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if err := e.Cars[car].Down(); err != nil {
		return err
	}

	e.Logging("D", car)
	return nil
}

func (e *Elevators) Stop(car int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if err := e.Cars[car].Stop(); err != nil {
		return err
	}

	e.Logging("S", car)
	return nil
}

func (e *Elevators) Open(car int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if err := e.Cars[car].Open(); err != nil {
		return err
	}

	e.Logging("O", car)
	return nil
}

func (e *Elevators) Close(car int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if err := e.Cars[car].Close(); err != nil {
		return err
	}

	e.Logging("C", car)
	return nil
}

func (e *Elevators) Enter(car int, users []int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if e.Cars[car].Status != "OPENED" {
		return fmt.Errorf("Wrong action")
	}

	if e.Cars[car].IsFull() {
		return fmt.Errorf("Exceeding the capacity of the elevator")
	}

	for _, id := range users {
		var found int = -1
		for i, v := range e.Calls {
			if (v.Timestamp <= e.Timestamp && v.Id == id) {
				found = i
				break
			}
		}

		if found < 0 {
			return fmt.Errorf("No passenger for %d", id)
		}

		c := e.Calls[found]
		e.Calls = append(e.Calls[:found], e.Calls[found+1:]...)

		if err := e.Cars[car].Enter(c); err != nil {
			e.Calls = append(e.Calls, c)
			return err
		}

		e.Outputs = append(e.Outputs, mydb.Output{e.Timestamp, id, "enter", e.Cars[car].Floor, car})

		log := fmt.Sprintf("o,%d,%d,enter,%d", e.Timestamp, id, e.Cars[car].Floor)
		e.Logs = append(e.Logs, log)
	}

	e.Logging("E", car)

	return nil
}

func (e *Elevators) Exit(car int, users []int) error {
	if car >= len(e.Cars) {
		return fmt.Errorf("Wrong car number")
	}

	if e.Cars[car].Status != "OPENED" {
		return fmt.Errorf("Wrong action")
	}

	for _, id := range users {
		c, err := e.Cars[car].Exit(id)
		if err != nil {
			return err
		}

		if c.End != e.Cars[car].Floor {
			e.Calls = append(e.Calls, mydb.Call{c.Id, e.Timestamp, e.Cars[car].Floor, c.End})
		}

		e.Outputs = append(e.Outputs, mydb.Output{e.Timestamp, id, "exit", e.Cars[car].Floor, car})

		log := fmt.Sprintf("o,%d,%d,exit,%d", e.Timestamp, id, e.Cars[car].Floor)
		e.Logs = append(e.Logs, log)
	}

	e.Logging("X", car)

	return nil
}

func (e *Elevators) GetOnCalls() (mydb.Calls, int) {
	call := make(mydb.Calls, 0)

	for _, v := range e.Calls {
		if (v.Timestamp <= e.Timestamp) {
			call = append(call, v)

			if v.Timestamp > e.LastCallTs {
				log := fmt.Sprintf("i,%d,%d,%d,%d", v.Timestamp, v.Id, v.Start, v.End)
				e.Logs = append(e.Logs, log)
			}
		}
	}

	e.LastCallTs = e.Timestamp

	return call, len(e.Calls) - len(call)
}
