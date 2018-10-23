package api

import (
    "fmt"
    "sort"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
    "2019-blind-2nd-elevator/elevator/internal/mydb"
    el "2019-blind-2nd-elevator/elevator/internal/elevator"
)

var (
	elevators *el.Elevators
)

func Start(w http.ResponseWriter, r *http.Request) {
	var err error
    params := mux.Vars(r)

	if params["userid"] == "" {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
	}

	problem, count := 0, 1

    if problem, err = strconv.Atoi(params["problem"]); err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    if count, err = strconv.Atoi(params["count"]); err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    if problem < 0 || problem > 2 || count < 1 || count > 4 {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    if elevators, err = el.NewElevators(problem, count); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return
    }

	isEnd := elevators.IsEnd()
	if isEnd {
		elevators.Done()
	}

    c := struct {
        Token string `json:"token"`
        Timestamp int `json:"timestamp"`
        Cars []el.Elevator `json:"elevators"`
        IsEnd bool `json:"is_end"`
    }{elevators.Token, elevators.Timestamp, elevators.Cars, isEnd}

    writeJson(w, c)
}

func OnCalls(w http.ResponseWriter, r *http.Request) {
	if err := checkToken(w, r); err != nil {
		return
	}
    
    call, _ := elevators.GetOnCalls()

	isEnd := elevators.IsEnd()
	if isEnd {
		elevators.Done()
	}

    c := struct {
        Token string `json:"token"`
        Timestamp int `json:"timestamp"`
        Cars []el.Elevator `json:"elevators"`
        mydb.Calls `json:"calls"`
        IsEnd bool `json:"is_end"`
    }{elevators.Token, elevators.Timestamp, elevators.Cars, call, isEnd}

    writeJson(w, c)
}

func Action(w http.ResponseWriter, r *http.Request) {
	if err := checkToken(w, r); err != nil {
		return
	}

    decoder := json.NewDecoder(r.Body)

	var input struct {
		Cmds []struct {
			CarId int `json:"elevator_id"`
			Cmd string `json:"command"`
			CallIds []int `json:"call_ids"`
		} `json:"commands"`
	}

    if err := decoder.Decode(&input); err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    if len(input.Cmds) != elevators.Count() {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    sort.Slice(input.Cmds, func(i, j int) bool {
        return input.Cmds[i].CarId < input.Cmds[j].CarId
    })


	for i, cmd := range input.Cmds {
		if i != cmd.CarId {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	// Order by command desc. to execute 'exit' before 'enter'.
	sort.Slice(input.Cmds, func(i, j int) bool {
		return input.Cmds[i].Cmd > input.Cmds[j].Cmd
	})

	for _, cmd := range input.Cmds {
		var err error
		switch cmd.Cmd {
		case "UP":
			err = elevators.Up(cmd.CarId)
		case "DOWN":
			err = elevators.Down(cmd.CarId)
		case "STOP":
			err = elevators.Stop(cmd.CarId)
		case "OPEN":
			err = elevators.Open(cmd.CarId)
		case "CLOSE":
			err = elevators.Close(cmd.CarId)
		case "ENTER":
			err = elevators.Enter(cmd.CarId, cmd.CallIds)
		case "EXIT":
			err = elevators.Exit(cmd.CarId, cmd.CallIds)
		default:
			err = fmt.Errorf("Not supported action: %s", cmd.Cmd)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	elevators.Tick()

	isEnd := elevators.IsEnd()
	if isEnd {
		elevators.Done()
	}

    c := struct {
        Token string `json:"token"`
        Timestamp int `json:"timestamp"`
        Cars []el.Elevator `json:"elevators"`
        IsEnd bool `json:"is_end"`
    }{elevators.Token, elevators.Timestamp, elevators.Cars, isEnd}

    writeJson(w, c)
}

func checkToken(w http.ResponseWriter, r *http.Request) error {
	elevatorId := r.Header.Get("X-Auth-Token")
	if elevatorId != elevators.Token {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("Wrong token")
	}
	return nil
}

func writeJson(w http.ResponseWriter, v interface{}) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(v)
}
