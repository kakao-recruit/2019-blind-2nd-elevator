package eval

import (
	"math"
	"fmt"
	"2019-blind-2nd-elevator/elevator/internal/mydb"
	log "github.com/sirupsen/logrus"
)


func haveIncorrectData(inputs mydb.Calls, outputs mydb.Outputs) error {
	for _, o := range outputs {
		if o.Timestamp < 0 {
			log.Debugf("Timestamp is negative")
			return fmt.Errorf("Incorrect outputs")
		}

		if  o.Id < 0 || o.Id >= len(inputs) {
			log.Debugf("wrong id")
			return fmt.Errorf("Incorrect outputs")
		} 

		if o.Floor < 0 || o.Floor > 100 {
			log.Debugf("wrong floor")
			return fmt.Errorf("Incorrect outputs")
		}

		if o.Action != "enter" && o.Action != "exit" {
			log.Debugf("wrong action : %s", o.Action)
			return fmt.Errorf("Incorrect outputs")
		}
	}

	return nil
}

func doesAnyoneWrongCar(inputs mydb.Calls, outputs mydb.Outputs) error {
	cars := make([]int, len(inputs))

	for _, output := range outputs {
		if output.Action == "exit" {
			if cars[output.Id] != output.Car {
				return fmt.Errorf("Different Car")
			}
		} else {
			cars[output.Id] = output.Car
		}
	}

	return nil
}

func doesAnyoneMissing(inputs mydb.Calls, outputs mydb.Outputs) error {
	serviced := make([]bool, len(inputs))
	for i := range serviced {
		serviced[i] = false
	}

	for _, output := range outputs {
		if output.Action == "exit" {
			serviced[output.Id] = true
		}
	}

	for _, s := range serviced {
		if s == false {
			return fmt.Errorf("Not all user has moved")
		}
	}

	return nil
}

func doesAnyoneMovingTooFast(inputs mydb.Calls, outputs mydb.Outputs) error {

	lastTs := []int{0, 0, 0, 0}
	lastFloor := []int{1, 1, 1, 1}

	for _, o := range outputs {
		car := o.Car

		move := int(math.Abs(float64(o.Floor - lastFloor[car])))
		ellapse := o.Timestamp - lastTs[car]

		if move > ellapse {
			return fmt.Errorf("Incorrect movements")
		}

		lastTs[car], lastFloor[car] = o.Timestamp, o.Floor
	}

	return nil
}

func doesAnyoneWrongDestination(inputs mydb.Calls, outputs mydb.Outputs) error {
	firstEnter := make([]int, len(inputs))
	lastExit := make([]int, len(inputs))

	for i := range inputs {
		firstEnter[i] = -1
		lastExit[i] = -1
	}

	for _, o := range outputs {
		if o.Action == "enter" && firstEnter[o.Id] < 0 {
			firstEnter[o.Id] = o.Floor
		}
		if o.Action == "exit" {
			lastExit[o.Id] = o.Floor
		}
	}

	for _, i := range inputs {
		if i.Start != firstEnter[i.Id] {
			return fmt.Errorf("Impossible movement")
		}
		if i.End != lastExit[i.Id] {
			return fmt.Errorf("Wrong destination")
		}
	}

	return nil

}

func doesAnyoneMismatchEnterExit(inputs mydb.Calls, outputs mydb.Outputs) error {
	enter := make([]int, len(inputs))
	for _, o := range outputs {
		if o.Action == "enter" {
			enter[o.Id] += 1
		} else if o.Action == "exit" {
			enter[o.Id] -= 1
		}

		if enter[o.Id] > 1 {
			return fmt.Errorf("Impossible movement")
		}
	}

	for _, e := range enter {
		if e != 0 {
			return fmt.Errorf("Not all user has moved")
		}
	}

	return nil
}

func doesAnyoneExitBeforeEnter(inputs mydb.Calls, outputs mydb.Outputs) error {
	enter := make([]int, len(inputs))
	for i := range enter {
		enter[i] = 99999999
	}

	for _, o := range outputs {
		if o.Action == "enter" {
			enter[o.Id] = o.Timestamp
		} else if o.Action == "exit" {
			if o.Timestamp < enter[o.Id] {
				return fmt.Errorf("Impossible movement")
			}
		}
	}

	return nil
}

func AverageWaitTime(inputs mydb.Calls, outputs mydb.Outputs) float64 {
	req := make([]int, len(inputs))
	for i := range req {
		req[i] = inputs[i].Timestamp
	}

	tot, cnt := 0.0, 0.0

	for _, o := range outputs {
		if o.Action == "enter" {
			tot += float64(o.Timestamp - req[o.Id])
			cnt += 1
		} else if o.Action == "exit" {
			// to support transfer case
			req[o.Id] = o.Timestamp
		}
	}

	if cnt == 0 {
		return 0.0
	}

	return tot / cnt
}

func AverageTravelTime(inputs mydb.Calls, outputs mydb.Outputs) float64 {
	enter := make([]int, len(inputs))

	tot, cnt := 0.0, 0.0
	for _, o := range outputs {
		if o.Action == "enter" {
			enter[o.Id] = o.Timestamp
		} else if o.Action == "exit" {
			tot += float64(o.Timestamp - enter[o.Id])
			cnt += 1
		}
	}

	if cnt == 0.0 {
		return 0.0
	}

	return tot / cnt
}

func AverageTotalTime(inputs mydb.Calls, outputs mydb.Outputs) float64 {
	lastExit := make([]int, len(inputs))

	for _, o := range outputs {
		if o.Action == "exit" {
			lastExit[o.Id] = o.Timestamp
		}
	}

	tot, cnt := 0.0, 0.0
	for i, input := range inputs {
		tot += float64(lastExit[i] - input.Timestamp)
		cnt += 1.0
	}

	if cnt == 0.0 {
		return 0.0
	}

	return tot / cnt
}

func IsValid(inputs mydb.Calls, outputs mydb.Outputs) error {
	funcs := [](func(mydb.Calls, mydb.Outputs) (error)) {
		haveIncorrectData,
		doesAnyoneWrongCar,
		doesAnyoneMissing,
		doesAnyoneWrongDestination,
		doesAnyoneMismatchEnterExit,
		doesAnyoneExitBeforeEnter,
		doesAnyoneMovingTooFast,
	}

	for _, f := range funcs {
		if err := f(inputs, outputs); err != nil {
			return err
		}
	}

	return nil
}

func Evaluate(inputs mydb.Calls, outputs mydb.Outputs) (wait, travel, total float64, lastTs int, err error) {
	if err = IsValid(inputs, outputs); err != nil {
		return
	}

	wait = AverageWaitTime(inputs, outputs)
	travel = AverageTravelTime(inputs, outputs)
	total = AverageTotalTime(inputs, outputs)
	lastTs = outputs[len(outputs) - 1].Timestamp
	return
}
