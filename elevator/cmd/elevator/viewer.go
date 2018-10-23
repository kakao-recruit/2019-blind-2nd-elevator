package main

import (
	"os"
	"path"
	"path/filepath"
	"io/ioutil"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"

	"github.com/spf13/viper"
	"github.com/gorilla/mux"

	_ "2019-blind-2nd-elevator/elevator/config"
)



func Index(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.WriteString("<!doctype html><style>td { padding-right: 1em }</style><table><tr><th>Token")

	logDir := viper.GetString("LogDir")

	files, _ := ioutil.ReadDir(logDir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".log" {
			token := path.Base(file.Name())
			token = token[:len(token) - 4]
			buf.WriteString(fmt.Sprintf("<tr><td><a href=/viewer/trials/%s>%s</a>", token, token))
		}
	}

	http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(buf.Bytes()))
}

func Trials(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	type Call struct {
		ID   int `json:"id"`
		From int `json:"from"`
		To   int `json:"to"`
	}

	type Elevator struct {
		Floor      int    `json:"floor"`
		State      string `json:"state"`
		Command    string `json:"command"`
		Passengers []Call `json:"passengers"`
	}

	type State struct {
		Elevators []Elevator `json:"elevators"`
		Calls     []Call     `json:"calls"`
	}

	timeline := make([]State, 0)
	var calls []Call
	passengers := [][]Call{{}, {}, {}, {}}
	boarding := make(map[int]Call)

	logDir := viper.GetString("LogDir")
	path := filepath.Join(logDir, fmt.Sprintf("%s.log", token))

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer file.Close()

	var tmp, status string
	var wait, travel, total float64
	var lastTs int

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%s %s", &tmp, &status)
	}

	if status == "OK" && scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%f %f %f %d", &wait, &travel, &total, &lastTs)
	}

	for scanner.Scan() {
		t := strings.Split(scanner.Text(), ",")

		ts := atoi(t[1])
		for ts >= len(timeline) {
			timeline = append(timeline, State{
				Elevators: []Elevator{
					{0, "", "", []Call{}},
					{0, "", "", []Call{}},
					{0, "", "", []Call{}},
					{0, "", "", []Call{}},
				},
				Calls: append([]Call{}, calls...),
			})
		}

		m := &timeline[ts]
		switch t[0] {
		case "i":
			calls = append(calls, Call{atoi(t[2]), atoi(t[3]), atoi(t[4])})

		case "e":
			id, floor := atoi(t[2]), atoi(t[3])

			p := make([]int, 0)
			if t[6] != "" {
				for _, e := range strings.Split(t[6], ":") {
					p = append(p, atoi(e))
				}
			}

			P1:
			for _, c := range passengers[id] {
				for i, e := range p {
					if e == c.To {
						p[i] = -1
						continue P1
					}
				}
				panic(c)
			}

			if t[4] == "E" {
				for i, c := range boarding {
					if floor == c.From {
						for j, n := range p {
							if n == c.To {
								p[j] = -1
								passengers[id] = append(passengers[id], c)
								delete(boarding, i)
								break
							}
						}
					}
				}
			}

			m.Elevators[id] = Elevator{
				Floor:      floor,
				State:      t[5],
				Command:    t[4],
				Passengers: append([]Call{}, passengers[id]...),
			}

		case "o":
			id := atoi(t[2])
			switch t[3] {
			case "enter":
				for i, c := range calls {
					if id == c.ID {
						boarding[c.ID] = c
						calls = append(calls[0:i], calls[i+1:]...)
						m.Calls = append([]Call{}, calls...)
						break
					}
				}
			case "exit":
				P:
				for i, p := range passengers {
					for j, c := range p {
						if id == c.ID {
							p = append(p[0:j], p[j+1:]...)
							passengers[i] = p
							break P
						}
					}
				}
			default:
				panic(t[3])
			}
		}
	}

	body, err := json.MarshalIndent(timeline, "", "  ")
	if err != nil {
		panic(err)
	}

	html := fmt.Sprintf(htmlBody, wait, travel, total, lastTs, status, body)
	http.ServeContent(w, r, "trial.html", time.Time{}, strings.NewReader(html))
}

func atoi(a string) int {
	n, _ := strconv.Atoi(a)
	return n
}
