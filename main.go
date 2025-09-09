package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gocmd <add|list|read|modify|delete> [flags]")
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		fs := flag.NewFlagSet("add", flag.ExitOnError)
		title := fs.String("title", "", "title")
		desc := fs.String("desc", "", "description")
		fs.Parse(os.Args[2:])
		// fmt.Printf("Adding task: title=%q desc=%q\n", *title, *desc)

		// Store it to the json file
		s, err := loadStore("./tasks.json")
		if err != nil {
			fmt.Println("Error loading store:", err.Error())
			return
		}
		task := Task{
			Title:       *title,
			Description: *desc,
		}
		err = s.save(task)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	case "modify":
		fs := flag.NewFlagSet("modify", flag.ExitOnError)
		id := fs.Int("id", 0, "id")

		title := fs.String("title", "", "Usage: --title \"title of the task\"")
		desc := fs.String("desc", "", "Usage: --desc \"description of the task\"")
		isCompleted := fs.Bool("completed", false, "Usage: --completed true")
		fs.Parse(os.Args[2:])

		if *id == 0 {
			fmt.Println("please provide valid --id")
			return
		}

		s, err := loadStore("./tasks.json")
		if err != nil {
			fmt.Println("Error loading store: ", err.Error())
			return
		}

		// Find existing task
		idx, err := s.GetTask(*id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		orig := s.Tasks[idx]

		// Update only when flag actually provided
		// Use Visit to see which flags were set by user
		updated := orig // start with original
		fs.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "title":
				updated.Title = *title
			case "desc":
				updated.Description = *desc
			case "completed":
				updated.Done = *isCompleted
			}
		})

		// Now persist update
		if err := s.update(updated, *id); err != nil {
			fmt.Println("update failed:", err)
			return
		}

		fmt.Println("MSG: The task has been successfully updated")
		updated.Print()
		return

	case "list":
		s, err := loadStore("./tasks.json")
		if err != nil {
			fmt.Println("Error loading store:", err)
			return
		}
		s.list()
	case "read":
		fs := flag.NewFlagSet("read", flag.ExitOnError)
		id := fs.Int("id", 0, "Usage: --id n ; where n is integer")
		fs.Parse(os.Args[2:])
		s, err := loadStore("./tasks.json")
		if err != nil {
			fmt.Println("Error loading store:", err)
			return
		}
		task, err := s.search(*id)
		if err != nil {
			fmt.Println(err.Error())
			return // avoid nil deref
		}
		task.Print()
	case "delete":
		fs := flag.NewFlagSet("delete", flag.ExitOnError)
		id := fs.Int("id", 0, "id")
		fs.Parse(os.Args[2:])
		s, err := loadStore("./tasks.json")
		if err != nil {
			fmt.Println("Error loading store:", err)
			return
		}
		if err := s.delete(*id); err != nil {
			fmt.Println("delete failed:", err)
			return
		}
		fmt.Println("Successfully deleted task")
	default:
		fmt.Println("Unknown action")
	}

}
