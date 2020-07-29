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
	errr     error
	question string
	answer   string
	correct  uint
	problems uint
}

func main() {

	curr := status{}
	curr.fileName = flag.String("csv", "problems.csv", "CSV file in 'question, answer' format")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*curr.fileName)
	curr.state, curr.errr = "open", err
	handleError(&curr)

	r := csv.NewReader(file)
	r.ReuseRecord, r.TrimLeadingSpace = true, true
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for {
		line, err := r.Read()
		curr.state, curr.errr = "read", err
		handleError(&curr)
		curr.problems++

		curr.question, curr.answer = line[0], line[1]
		fmt.Printf("Problem #%d: %s = ", curr.problems, curr.question)
		answerCh := make(chan string)

		go func() {
			var submission string
			fmt.Scanf("%s\n", &submission)
			answerCh <- submission
		}()

		select {
		case <-timer.C:
			exit(fmt.Sprintf("\nTime expired"))
		case submission := <-answerCh:
			if curr.answer == submission {
				curr.correct++
			}
		}
	}
}

func handleError(curr *status) {
	switch curr.state {
	case "open":
		if curr.errr != nil {
			exit(fmt.Sprintf("Failed to open %s", *curr.fileName))
		}
	case "read":
		if curr.errr == io.EOF {
			exit(fmt.Sprintf("You scored %d out of %d", curr.correct, curr.problems))
		}
		if curr.errr != nil {
			exit(fmt.Sprintf("Failed to read CSV %v", curr.errr))
		}
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
