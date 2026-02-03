package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

func generateId() string {
	var id [3]byte
	_, err := rand.Read(id[:])
	if err != nil {
		log.Fatal("Error reading random bytes", err)
	}
	return fmt.Sprintf("%x", id)
}

type ToDo struct {
	ID         string
	Text       string
	IsComplete bool
	Date       time.Time
}

func (t *ToDo) Create() {
}

func (t *ToDo) Read() {
}

func (t *ToDo) Update() {
}

func (t *ToDo) Delete() {
}

// ReadFile
// WriteToFile
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

	val := generateId()

	fmt.Printf("Random Value: %v\n", val)
	fmt.Printf("This is my To Do program\n")

	/// maybe I can have a switch function that checks what has flag has been
	/// called, probably put it in some function
}
