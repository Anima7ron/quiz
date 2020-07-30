package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type status struct {
	state    string
	fileName *string
	fault    error
	problem  string
	solution string
	score    uint
	lines    uint
}

func main() {

	my := status{}
	my.fileName = flag.String("csv", "problems.csv", "CSV file in 'problem, solution' format")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*my.fileName)
	my.state, my.fault = "open", err
	handleError(&my)

	r := csv.NewReader(file)
	r.ReuseRecord, r.TrimLeadingSpace = true, true
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for {
		line, err := r.Read()
		my.state, my.fault = "read", err
		handleError(&my)
		my.lines++

		my.problem, my.solution = line[0], line[1]
		fmt.Printf("Problem #%d: %s = ", my.lines, my.problem)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			exit(fmt.Sprintf("\nTime expired"))
		case answer := <-answerCh:
			if my.solution == answer {
				my.score++
			}
		}
	}
}

func handleError(my *status) {
	switch my.state {
	case "open":
		if my.fault != nil {
			exit(fmt.Sprintf("Failed to open %s", *my.fileName))
		}
	case "read":
		if my.fault == io.EOF {
			exit(fmt.Sprintf("You scored %d out of %d", my.score, my.lines))
		}
		if my.fault != nil {
			exit(fmt.Sprintf("Failed to read CSV %v", my.fault))
		}
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
