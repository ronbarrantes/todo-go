package main

import (
	"crypto/rand"
	"flag"
	"fmt"
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

type Store struct {
	db *gorm.DB
}

func (store *Store) GetDatabase() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(home, ".local", "share", "todo-go")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return err
	}

	gormPath := filepath.Join(appDir, "gorm.db")

	db, err := gorm.Open(sqlite.Open(gormPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	store.db = db
	return nil
}

func generateId() (string, error) {
	var id [6]byte
	_, err := rand.Read(id[:])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", id), nil
}

// CREATE A FIND FUNCTION FOR SHORTENED IDS
func printToDo(t *ToDo) {
	completed := " "
	if t.IsCompleted {
		completed = "x"
	}

	fmt.Printf("- [%s] (%s) %s\n", completed, t.ID, t.Text)
}

func (store *Store) Create(t string) error {
	id, err := generateId()
	if err != nil {
		return err
	}

	item := &ToDo{
		ID:          id,
		Text:        t,
		IsCompleted: false,
	}

	if err := store.db.Create(item).Error; err != nil {
		fmt.Printf("To-do %s not created\n", id)
		return err
	}

	fmt.Printf("To-do %s created\n", id)
	return nil
}

func (store *Store) Read() error {
	var todos []*ToDo
	if err := store.db.Find(&todos).Error; err != nil {
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

func (store *Store) Update(td *ToDo) error {
	if err := store.db.Model(&td).Update("Text", td.Text).Error; err != nil {
		fmt.Printf("To-do %s not updated", td.ID)
		return err
	}
	fmt.Printf("To-do %s updated", td.ID)
	return nil
}

func (store *Store) Delete(id string) error {
	if err := store.db.Delete(&ToDo{ID: id}).Error; err != nil {
		fmt.Printf("To-do %s not deleted", id)
		return err
	}
	fmt.Printf("To-do %s deleted", id)
	return nil
}

func (store *Store) ToggleTodo(id string) error {
	var todo ToDo
	if err := store.db.Where("id = ?", id).First(&todo).Error; err != nil {
		return err
	}

	completed := !todo.IsCompleted

	if err := store.db.Model(&todo).Update("is_completed", completed).Error; err != nil {
		return err
	}
	return nil
}

func main() {
	s := time.Now()
	var store Store
	err := store.GetDatabase()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	store.db.AutoMigrate(&ToDo{})

	defer func() {
		duration := time.Since(s)
		fmt.Printf("\nThis program took %v to run\n", duration)
	}()

	createFlag := flag.String("c", "", "Create a to do")
	readFlag := flag.Bool("r", false, "Read all to todos")
	// findFlag := flag.String("f", "", "Find a to todo")
	updateFlag := flag.String("u", "", "Update a to do")
	textFlag := flag.String("t", "", "Update the text of a to do")
	toggleCompleteFlag := flag.String("x", "", "Update completed state")
	deleteFlag := flag.String("d", "", "Delete to do")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	todo := ToDo{}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		fmt.Println("Listing todos by default...")
		err := store.Read()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		return
	}

	switch {
	case *readFlag:
		fmt.Println("Reading ....")
		err := store.Read()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *createFlag != "":
		err := store.Create(*createFlag)
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
		err := store.Update(&todo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case *toggleCompleteFlag != "":
		err := store.ToggleTodo(*toggleCompleteFlag)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case *deleteFlag != "":
		err := store.Delete(*deleteFlag)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	default:
		fmt.Println("Unknown flag or no action")
		flag.Usage()
		os.Exit(1)
	}
}
