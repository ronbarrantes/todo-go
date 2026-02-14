package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	db "ronb.co/todo-go/store"
)

func run() (*db.Store, error) {
	var store db.Store
	err := store.GetDatabase()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	store.DB.AutoMigrate(&db.ToDo{})

	return &store, nil
}

type flags struct {
	create string
	list   bool
	update string
	text   string
	toggle string
	delete string
	help   bool
}

func parseFlags() flags {
	createFlag := flag.String("c", "", "Create a to-do")
	listFlag := flag.Bool("l", false, "List all to-dos")
	// findFlag := flag.String("f", "", "Find a to todo")
	updateFlag := flag.String("u", "", "Update a to-do")
	textFlag := flag.String("t", "", "Update the text of a to-do")
	toggleCompleteFlag := flag.String("x", "", "Toggle to-do completion")
	deleteFlag := flag.String("d", "", "Delete to-do")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	return flags{
		create: *createFlag,
		list:   *listFlag,
		update: *updateFlag,
		text:   *textFlag,
		toggle: *toggleCompleteFlag,
		delete: *deleteFlag,
		help:   *helpFlag,
	}
}

func main() {
	secs := time.Now()
	defer func() {
		duration := time.Since(secs)
		fmt.Printf("-----\nThis program took %v to run\n", duration)
	}()

	store, err := run()
	f := parseFlags()
	if err != nil {
		fmt.Printf("error, %v", err)
	}

	todo := db.ToDo{}

	if f.help {
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
	case f.list:
		fmt.Println("Reading ....")
		err := store.Read()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case f.create != "":
		err := store.Create(f.create)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case f.update != "":
		if f.text == "" {
			fmt.Println("please provide the text with -t")
			os.Exit(1)
		}

		todo.ID = f.update
		todo.Text = f.text
		err := store.Update(&todo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case f.toggle != "":
		err := store.ToggleTodo(f.toggle)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	case f.delete != "":
		err := store.Delete(f.delete)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

	default:
		fmt.Println("Unknown flag or no action")
		flag.Usage()
		os.Exit(1)
	}
}
