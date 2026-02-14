package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	db "ronb.co/todo-go/store"
)

func main() {
	s := time.Now()
	var store db.Store
	err := store.GetDatabase()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	store.DB.AutoMigrate(&db.ToDo{})

	defer func() {
		duration := time.Since(s)
		fmt.Printf("-----\nThis program took %v to run\n", duration)
	}()

	createFlag := flag.String("c", "", "Create a to-do")
	readFlag := flag.Bool("l", false, "List all to-dos")
	// findFlag := flag.String("f", "", "Find a to todo")
	updateFlag := flag.String("u", "", "Update a to-do")
	textFlag := flag.String("t", "", "Update the text of a to-do")
	toggleCompleteFlag := flag.String("x", "", "Toggle to-do completion")
	deleteFlag := flag.String("d", "", "Delete to-do")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	todo := db.ToDo{}

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
