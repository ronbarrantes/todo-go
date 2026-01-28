package main

import (
	"fmt"
	"time"
)

// going to make a todo
// FULL CRUD
// save items to a json
// do it via cli with flags
// -a : --add
// -d : --done
// -l : --list
// -d: --delete

func main() {
	s := time.Now()

	defer func() {
		duration := time.Since(s)
		fmt.Printf("This program took %v to complete\n", duration)
	}()

	fmt.Printf("This is my To Do program\n")
}
