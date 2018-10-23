package mydb

import (
	"os"
	"fmt"
	"bufio"
	log "github.com/sirupsen/logrus"
	"path"

	"github.com/spf13/viper"
)

type Call struct {
	Id        int `json:"id"`
	Timestamp int `json:"timestamp"`
	Start     int `json:"start"`
	End       int `json:"end"`
}

type Calls []Call

type Output struct {
	Timestamp int
	Id        int
	Action    string
	Floor     int
	Car       int
}

type Outputs []Output

var (
	problems []Calls
)

func read(path string) Calls {
	var inputs Calls
	var ts, uid, start, end int

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err, path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d,%d,%d", &ts, &uid, &start, &end)
		inputs = append(inputs, Call{uid, ts, start, end})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return inputs
}

func init() {
	dir := viper.GetString("DatasetDir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal(err, dir)
	}

	logDir := viper.GetString("LogDir")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		log.Fatal(err, logDir)
	}

	problems = make([]Calls, 0, 3)
	problems = append(problems, read(path.Join(dir, "p0.in")))
	problems = append(problems, read(path.Join(dir, "p1.in")))
	problems = append(problems, read(path.Join(dir, "p2.in")))
}

func GetCalls(problem int) (Calls) {
	return problems[problem]
}
