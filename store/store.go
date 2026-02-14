package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	u "ronb.co/todo-go/utils"
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
	DB *gorm.DB
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

	store.DB = db
	return nil
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
	id, err := u.GenerateId()
	if err != nil {
		return err
	}

	item := &ToDo{
		ID:          id,
		Text:        t,
		IsCompleted: false,
	}

	if err := store.DB.Create(item).Error; err != nil {
		fmt.Printf("To-do %s not created\n", id)
		return err
	}

	fmt.Printf("To-do %s created\n", id)
	return nil
}

func (store *Store) Read() error {
	var todos []*ToDo
	if err := store.DB.Find(&todos).Error; err != nil {
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

func (store *Store) FindShortId(prefix string) (*ToDo, error) {
	var todos []ToDo

	if len(prefix) < 3 {
		return nil, fmt.Errorf("Prefix needs to be longer than 3 characters")
	}

	if err := store.DB.Where("id LIKE ?", prefix+"%").Find(&todos).Error; err != nil {
		return nil, err
	}

	switch len(todos) {
	case 0:
		return nil, fmt.Errorf("No to-do found")

	case 1:
		return &todos[0], nil

	default:
		return nil, fmt.Errorf("Too ambiguous")
	}
}

func (store *Store) Update(td *ToDo) error {
	todo, err := store.FindShortId(td.ID)
	if err != nil {
		return err
	}

	if err := store.DB.Model(todo).Update("Text", td.Text).Error; err != nil {
		fmt.Printf("To-do %s not updated\n", td.ID)
		return err
	}

	fmt.Printf("To-do %s updated\n", td.ID)
	return nil
}

func (store *Store) Delete(id string) error {
	todo, err := store.FindShortId(id)
	if err != nil {
		return err
	}
	if err := store.DB.Delete(todo).Error; err != nil {
		fmt.Printf("To-do %s not deleted\n", id)
		return err
	}
	fmt.Printf("To-do %s deleted\n", id)
	return nil
}

func (store *Store) ToggleTodo(id string) error {
	todo, err := store.FindShortId(id)
	if err != nil {
		return err
	}
	completed := !todo.IsCompleted
	if err := store.DB.Model(todo).Update("is_completed", completed).Error; err != nil {
		return err
	}

	fmt.Printf("To-do %s toggled\n", id)
	return nil
}
