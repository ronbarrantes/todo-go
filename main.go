package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type ToDo struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	IsCompleted bool      `json:"is_completed"`
	Date        time.Time `json:"date"`
}

type ToDoUpdate struct {
	Text        *string
	IsCompleted *bool
}

func getUserDataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(home, ".local", "share", "todo-go")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "todo.json"), nil
}

func generateId() string {
	var id [6]byte
	_, err := rand.Read(id[:])
	if err != nil {
		log.Fatal("Error reading random bytes", err)
	}
	return fmt.Sprintf("%x", id)
}

// CREATE A FIND FUNCTION FOR SHORTENED IDS

func printToDo(t *ToDo) {
	completed := " "
	if t.IsCompleted {
		completed = "x"
	}

	fmt.Printf("- [%s] (%s) %s\n", completed, t.ID, t.Text)
}

func writeJSON(t []*ToDo) error {
	jsonBlob, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}

	path, err := getUserDataPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonBlob, 0644)
}

func readJSON() ([]*ToDo, error) {
	path, err := getUserDataPath()
	if err != nil {
		return nil, err
	}

	todos := []*ToDo{}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return todos, nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return todos, nil
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
		item := &ToDo{
			ID:          val,
			Date:        time.Now(),
			Text:        t.Text,
			IsCompleted: false,
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
		return fmt.Errorf("there are no to dos")
	}

	for _, todo := range jData {
		printToDo(todo)
	}

	return nil
}

func (t *ToDo) ApplyUpdate(u ToDoUpdate) {
	if u.Text != nil {
		t.Text = *u.Text
	}
	if u.IsCompleted != nil {
		t.IsCompleted = *u.IsCompleted
	}
}

func (t *ToDo) Update(u ToDoUpdate) error {
	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
		if len(td) == 0 {
			fmt.Println("nothing to do!")
			return []*ToDo{}, nil
		}

		for i, jToDo := range td {
			if jToDo.ID == t.ID {
				jToDo.ApplyUpdate(u)
				td[i] = jToDo
				return td, nil
			}
		}

		return nil, fmt.Errorf("to do %s not found", t.ID)
	})
}

func (t *ToDo) SetComplete() error {
	completed := true
	return t.Update(ToDoUpdate{IsCompleted: &completed})
}

// func (t *ToDo) ToggleTodo() error {
// 	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
// 		if len(td) == 0 {
// 			fmt.Println("nothing to do!")
// 			return []*ToDo{}, nil
// 		}

// 		for _, jToDo := range td {
// 			if jToDo.ID == t.ID {
// 				completed := true
// 				notCompleted := false
// 				if t.IsCompleted != nil || !*jToDo.IsCompleted {
// 					jToDo.IsCompleted = &completed
// 				} else {
// 					jToDo.IsCompleted = &notCompleted
// 				}

// 				return td, nil
// 			}
// 		}

// 		return nil, fmt.Errorf("to do %s not found", t.ID)
// 	})
// }

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

	defer func() {
		duration := time.Since(s)
		fmt.Printf("\nThis program took %v to run\n", duration)
	}()

	createFlag := flag.String("c", "", "Create a to do")
	readFlag := flag.Bool("r", false, "Read all to todos")
	// findFlag := flag.String("f", "", "Find a to todo")
	updateFlag := flag.String("u", "", "Update a to do")
	textFlag := flag.String("t", "", "Update a to do")
	completeFlag := flag.String("x", "", "Update a to do")
	deleteFlag := flag.String("d", "", "Delete to do")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	var upd ToDoUpdate
	todo := ToDo{}

	if *textFlag != "" && *updateFlag != "" {
		upd.Text = textFlag
	}

	if *completeFlag != "" {
		completed := true
		upd.IsCompleted = &completed
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		fmt.Println("Listing todos by default...")
		return
	}

	switch {
	case *readFlag:
		fmt.Println("Reading ....")
		err := todo.Read()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *createFlag != "":
		fmt.Printf("the flag is %s\n", *createFlag)
		todo.Text = *createFlag
		err := todo.Create()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *updateFlag != "":
		if *textFlag == "" {
			fmt.Printf(`please provide the text with -t`)
			os.Exit(1)
		}

		todo.ID = *updateFlag
		err := todo.Update(upd)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *completeFlag != "":
		err := todo.SetComplete()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *deleteFlag != "":
		fmt.Printf("the flag is %s\n", *createFlag)
		todo.ID = *deleteFlag
		err := todo.Delete()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	default:
		fmt.Println("Unknown flag or no action")
		flag.Usage()
		os.Exit(1)
	}
}
