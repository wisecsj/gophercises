package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFileName := flag.String("name", "problems.csv", "csv file name")
	timeLimit := flag.Int64("limit", 30, "time limit(s)")
	flag.Parse()

	filePtr, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("%s", err))
	}

	r := csv.NewReader(filePtr)
	records, err := r.ReadAll()
	if err != nil {
		exit("csv file parse failed!")
	}
	problems := parseRecords(records)

	notStart := true
	fmt.Print("Please press enter to start quiz...")
	var input string
	for notStart {
		fmt.Scanln(&input)
		// fmt.Print("-",input,"-")
		if input == "" {
			notStart = false
		}
	}

	ch := make(chan int)
	go func(limit int64, ch chan int) {
		var v int
		for {
			select {
			case v = <-ch:
			case <-time.After(time.Duration(limit) * time.Second):
				fmt.Println("\nTime limit")
				fmt.Printf("You scored %d out of %d.\n", v, len(problems))
				os.Exit(0)
			}
		}
	}(*timeLimit, ch)

	var correct int

	for i, p := range problems {
		fmt.Printf("%d.%s=", i+1, p.q)
		var input string
		fmt.Scanf("%s\n", &input)
		if input == p.a {
			correct++
			ch <- correct
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))

}
func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))
	for i, r := range records {
		problems[i] = problem{
			q: r[0],
			a: strings.TrimSpace(r[1]),
		}
	}
	return problems
}
func exit(s string) {
	fmt.Println(s)
	os.Exit(1)
}
