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
	"slices"
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
	completed := " "
	if *t.IsCompleted {
		completed = "x"
	}
	fmt.Printf("- [%s] id: %s | item %s\n", completed, t.ID, t.Text)
}

// func init() {
// }

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

func updateStore(fn func([]*ToDo) ([]*ToDo, error)) error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	if len(jData) == 0 {
		fmt.Println("nothing to do!")
		return nil
	}

	toReadAndWrite, err := fn(jData)
	if err != nil {
		return err
	}

	err = writeJSON(toReadAndWrite)
	if err != nil {
		return err
	}

	return nil
}

func (t *ToDo) Create() error {
	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
		val := generateId()

		completed := false
		item := &ToDo{
			ID:          val,
			Date:        time.Now(),
			Text:        t.Text,
			IsCompleted: &completed,
		}

		td = append(td, item)

		printToDo(item)
		return td, nil
	})
}

func (t *ToDo) Read() error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	if len(jData) == 0 {
		fmt.Println("nothing to do!")
	}

	for _, todo := range jData {
		printToDo(todo)
	}

	return nil
}

func (t *ToDo) Update() error {
	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
		for _, jToDo := range td {
			if jToDo.ID == t.ID {
				if t.Text != "" {
					jToDo.Text = t.Text
				}

				if t.IsCompleted != nil {
					jToDo.IsCompleted = t.IsCompleted
				}

				return td, nil
			}
		}

		return nil, fmt.Errorf("to do %s not found", t.ID)
	})
}

func (t *ToDo) ToggleTodo() error {
	jData, err := readJSON()
	if err != nil {
		return err
	}

	if len(jData) == 0 {
		fmt.Println("nothing to do!")
		return nil
	}

	found := false
	for i, jToDo := range jData {
		if jToDo.ID == t.ID {
			completed := true
			notCompleted := false
			if t.IsCompleted != nil || !*jToDo.IsCompleted {
				jToDo.IsCompleted = &completed
			} else {
				jToDo.IsCompleted = &notCompleted
			}

			found = true
			jData[i] = jToDo
			break
		}
	}

	if !found {
		return fmt.Errorf("to do %s not found", t.ID)
	}

	err = writeJSON(jData)
	if err != nil {
		return err
	}

	return nil
}

func (t *ToDo) SetComplete() error {
	completed := true
	t.IsCompleted = &completed
	return t.Update()
}

func (t *ToDo) Delete() error {
	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
		for i, todo := range td {
			if todo.ID == t.ID {
				return slices.Delete(td, i, i+1), nil
			}
		}

		return nil, fmt.Errorf("to do with id %s not found", t.ID)
	})
}

func main() {
	s := time.Now()
	completed := false
	item := &ToDo{
		ID:          "baef51",
		IsCompleted: &completed,
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

	err := item.Delete()
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	//	fmt.Printf("Random Value: %v\n", val)

	/// maybe I can have a switch function that checks what has flag has been
	/// called, probably put it in some function
}
