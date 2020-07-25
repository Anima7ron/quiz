package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	//Set flags for input CSV and time-limit
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	//Open CSV file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
		os.Exit(1)
	}
	//Parse CSV
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	problems := parseLines(lines)
	//Initialize timer and correct answer counter
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemloop:
	for i, p := range problems {
		//Present each problem
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		//Submit each answer
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		//Either timer runs out...
		case <-timer.C:
			fmt.Println()
			break problemloop
		//... Or answer is submitted
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	//Present Score
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

//Convert parsed problem to problem structure
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

//Struct problem makes quiz data portable
type problem struct {
	q string
	a string
}

//Game Over
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
