package main

import (
	"fmt"
	"time"
)

func statusUpdate() string { return "" }

func main() {
	c := time.Tick(1e9)
	for now := range c {
		fmt.Printf("%v %s\n", now, statusUpdate())
	}
}