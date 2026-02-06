package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
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
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	IsComplete bool      `json:"is_complete"`
	Date       time.Time `json:"date"`
}

func init() {
}

func writeJSON(t []*ToDo) error {
	jsonBlob, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}

	path := getUserPath()
	return os.WriteFile(path, jsonBlob, 0644)
}

func readJSON() ([]*ToDo, error) {
	path := getUserPath()
	jData := []*ToDo{}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return jData, nil
		}

		return nil, err
	}

	if err = json.Unmarshal(data, &jData); err != nil {
		return nil, err
	}

	return jData, nil
}

func (t *ToDo) Create() error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	updatedData := append(jData, t)

	fmt.Printf("id: %s | item %s", t.ID, t.Text)
	return writeJSON(updatedData)
}

func (t *ToDo) Read() {
	readJSON()
}

func (t *ToDo) Update() {
	// read

	// find

	// update

	// write
}

func (t *ToDo) Delete() {
	// read

	// find

	// remove

	// write
}

// FULL CRUD
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

	item := &ToDo{
		ID:         val,
		Date:       s,
		Text:       fmt.Sprintf("New item %s", val),
		IsComplete: false,
	}

	item.Create()

	//	item.Read()

	//	fmt.Printf("Random Value: %v\n", val)

	/// maybe I can have a switch function that checks what has flag has been
	/// called, probably put it in some function
}
