package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ToDo struct {
	ID          string `gorm:"primaryKey"`
	Text        string `gorm:"text"`
	IsCompleted bool   `gorm:"is_completed"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type DBContext struct {
	db *gorm.DB
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

func (ctx *DBContext) getDatabase() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(home, ".local", "share", "todo-go")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	ctx.db = db
	return nil
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

func (ctx *DBContext) Create(t *ToDo) error {
	id := generateId()
	item := &ToDo{
		ID:          id,
		Text:        t.Text,
		IsCompleted: false,
	}

	if err := ctx.db.Create(item).Error; err != nil {
		return err
	}

	fmt.Printf("To-do %s created\n", id)
	return nil
}

func (ctx *DBContext) Read() error {
	var todos []*ToDo
	if err := ctx.db.Find(&todos).Error; err != nil {
		return err
	}

	if len(todos) == 0 {
		return fmt.Errorf("There are no to-dos")
	}

	for _, todo := range todos {
		printToDo(todo)
	}

	return nil
}

func (ctx *DBContext) Update(td *ToDo) error {
	if err := ctx.db.Model(&td).Update("Text", &td.Text).Error; err != nil {
		// fmt.Errorf("to do %s not found", t.ID)
		return err
	}
	fmt.Printf("To-do %s updated", td.ID)
	return nil
}

func (ctx *DBContext) Delete(td *ToDo) error {
	if err := ctx.db.Delete(&td).Error; err != nil {
		// fmt.Errorf("to do %s not found", t.ID)
		return err
	}
	fmt.Printf("To-do %s deleted", td.ID)
	return nil
}

// func (t *ToDo) Update() error {
// 	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
// 		if len(td) == 0 {
// 			fmt.Println("nothing to do!")
// 			return []*ToDo{}, nil
// 		}

// 		todo, err := t.FindItem(td)
// 		if err != nil {
// 			return nil, err
// 		}

// 		todo.Text = t.Text
// 		return td, nil
// 	})
// }

// func (t *ToDo) FindItem(td []*ToDo) (*ToDo, error) {
// 	for _, todo := range td {
// 		if len(t.ID) >= 4 && s.HasPrefix(todo.ID, t.ID) {
// 			return todo, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("to do %s not found", t.ID)
// }

// func (t *ToDo) ToggleTodo() error {
// 	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
// 		if len(td) == 0 {
// 			fmt.Println("nothing to do!")
// 			return []*ToDo{}, nil
// 		}

// 		todo, err := t.FindItem(td)
// 		if err != nil {
// 			return nil, err
// 		}

// 		todo.IsCompleted = !todo.IsCompleted
// 		return td, nil
// 	})
// }

// func (t *ToDo) Delete() error {
// 	return updateStore(func(td []*ToDo) ([]*ToDo, error) {
// 		found, err := t.FindItem(td)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for i, todo := range td {
// 			if todo.ID == found.ID {
// 				return slices.Delete(td, i, i+1), nil
// 			}
// 		}

// 		return nil, fmt.Errorf("to do with id %s not found", t.ID)
// 	})
// }

func main() {
	s := time.Now()
	var dbCtx DBContext
	err := dbCtx.getDatabase()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	dbCtx.db.AutoMigrate(&ToDo{})

	defer func() {
		duration := time.Since(s)
		fmt.Printf("\nThis program took %v to run\n", duration)
	}()

	createFlag := flag.String("c", "", "Create a to do")
	readFlag := flag.Bool("r", false, "Read all to todos")
	// findFlag := flag.String("f", "", "Find a to todo")
	updateFlag := flag.String("u", "", "Update a to do")
	textFlag := flag.String("t", "", "Update the text of a to do")
	// toggleCompleteFlag := flag.String("x", "", "Update completed state")
	// deleteFlag := flag.String("d", "", "Delete to do")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	todo := ToDo{}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		fmt.Println("Listing todos by default...")
		err := dbCtx.Read()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		return
	}

	switch {
	case *readFlag:
		fmt.Println("Reading ....")
		err := dbCtx.Read()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *createFlag != "":
		todo.Text = *createFlag
		err := dbCtx.Create(&todo)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *updateFlag != "":
		if *textFlag == "" {
			fmt.Println("please provide the text with -t")
			os.Exit(1)
		}

		todo.ID = *updateFlag
		todo.Text = *textFlag
		err := dbCtx.Update(&todo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	// case *toggleCompleteFlag != "":
	// 	todo.ID = *toggleCompleteFlag
	// 	err := todo.ToggleTodo()
	// 	if err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 	}

	// case *deleteFlag != "":
	// 	todo.ID = *deleteFlag
	// 	err := todo.Delete()
	// 	if err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 	}

	default:
		fmt.Println("Unknown flag or no action")
		flag.Usage()
		os.Exit(1)
	}
}
