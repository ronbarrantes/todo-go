package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func getUserPath() string {
	currentUser, err := user.Current()
	check(err)
	return filepath.Join(currentUser.HomeDir, "Documents", "todo-go")
}

func generateId() string {
	var id [3]byte
	_, err := rand.Read(id[:])
	if err != nil {
		log.Fatal("Error reading random bytes", err)
	}
	return fmt.Sprintf("%x", id)
}

// check errors

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ToDo struct {
	ID         string
	Text       string
	IsComplete bool
	Date       time.Time
}

func writeJSON() {
	data := []byte("hello\ngo\n")
	path := getUserPath()
	err := os.WriteFile(path, data, 0644)
	check(err)

	// defer f.close()
}

func readJSON() {
}

func (t *ToDo) Create() {
	writeJSON()
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

	writeJSON()

	val := generateId()

	fmt.Printf("Random Value: %v\n", val)
	fmt.Printf("This is my To Do program\n")

	/// maybe I can have a switch function that checks what has flag has been
	/// called, probably put it in some function
}
