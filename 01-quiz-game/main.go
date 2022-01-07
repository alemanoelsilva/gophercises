package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const EXIT = 1

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit("Failed to open the CSV file")
	}

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		exit("Failed to part the provided CSV file.")
	}

	problems := parseLine(lines)

	correct := 0

	for index, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", index+1, problem.question)

		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
		answerChannel := make(chan string)

		go func() {
			var answer string
			// Scanf function is a blocking function.
			// we are using a goroutine to not block the timer into the select statement
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		// when the timer.C receive a message, the timer will be triggered
		case <-timer.C:
			showScore(correct, len(problems))
			return
		// when the answerChannel receive the answer set in the goroutine, the answer will be received
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}

	showScore(correct, len(problems))
}

type problem struct {
	question string
	answer   string
}

func parseLine(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for index, line := range lines {
		problems[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func showScore(corret int, total int) {
	fmt.Printf("\nYou scored %d out of %d\n", corret, total)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(EXIT)
}
