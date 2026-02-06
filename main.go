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

type ToDo struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	IsCompleted *bool     `json:"is_completed"`
	Date        time.Time `json:"date"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func printToDo(t *ToDo) {
	var completed string = " "
	if *t.IsCompleted {
		completed = "x"
	}
	fmt.Printf("- [%s] id: %s | item %s\n", completed, t.ID, t.Text)
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
	todos := []*ToDo{}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return todos, nil
		}

		return nil, err
	}

	if err = json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *ToDo) Create() error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	val := generateId()

	completed := false
	item := &ToDo{
		ID:          val,
		Date:        time.Now(),
		Text:        fmt.Sprintf("New item %s", val),
		IsCompleted: &completed,
	}

	updatedData := append(jData, item)

	fmt.Printf("id: %s | item %s\n", item.ID, item.Text)
	return writeJSON(updatedData)
}

// func (t *ToDo) Read

func (t *ToDo) Read() error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	if len(jData) == 0 {
		fmt.Println("Nothing to do!")
	}

	for _, todo := range jData {
		printToDo(todo)
	}

	return nil
}

func (t *ToDo) Update() error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	if len(jData) == 0 {
		fmt.Println("Nothing to do!")
		return nil
	}

	found := false
	for i, jToDo := range jData {
		if jToDo.ID == t.ID {
			if t.Text != "" {
				jToDo.Text = t.Text
			}

			if t.IsCompleted != nil {
				jToDo.IsCompleted = t.IsCompleted
			}

			fmt.Printf("Updated %v", jToDo)
			found = true
			jData[i] = jToDo
			break
		}
	}

	if !found {
		return fmt.Errorf("To do %s not found", t.ID)
	}

	err = writeJSON(jData)
	if err != nil {
		return err
	}

	return nil
}

func (t *ToDo) ToggleTodo() error {
	var isCompleted bool

	if t.IsCompleted != nil {
		isCompleted = !t.IsCompleted
	}

	todo := ToDo{
		ID:          t.ID,
		IsCompleted: &isCompleted,
	}

	return todo.Update()
}

func (t *ToDo) Delete() {
	// read

	// find

	// remove

	// write
}

func main() {
	s := time.Now()
	item := &ToDo{
		ID:   "dce04d",
		Text: "CDEF",
	}

	defer func() {
		duration := time.Since(s)
		fmt.Printf("\nThis program took %v to run\n", duration)
	}()

	// FULL CRUD
	// -a : --add
	// -d : --done
	// -l : --list
	// -d: --delete

	item.ToggleTodo()

	//	fmt.Printf("Random Value: %v\n", val)

	/// maybe I can have a switch function that checks what has flag has been
	/// called, probably put it in some function
}
