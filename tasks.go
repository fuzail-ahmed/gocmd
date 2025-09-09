package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (t *Task) Print() {
	fmt.Println("ID: ", t.ID)
	fmt.Println("Title: ", t.Title)
	fmt.Println("Description: ", t.Description)
	fmt.Println("Done: ", t.Done)
}

type Store struct {
	FilePath string
	Tasks    []Task
	nextID   int
}

// save adds task, persists and increments nextID
func (s *Store) save(task Task) error {
	task.ID = s.nextID
	s.nextID++
	s.Tasks = append(s.Tasks, task)

	data, err := json.MarshalIndent(s.Tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling tasks: %w", err)
	}

	if err := os.WriteFile(s.FilePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

// GetTask returns index in slice for a given ID, or error if not found.
func (s *Store) GetTask(ID int) (int, error) {
	for idx, task := range s.Tasks {
		if task.ID == ID {
			return idx, nil
		}
	}
	return -1, fmt.Errorf("task with id %d does not exist", ID)
}

// update updates only provided fields (caller must decide which fields to update)
func (s *Store) update(task Task, ID int) error {
	idx, err := s.GetTask(ID)
	if err != nil {
		return err
	}

	old := s.Tasks[idx]

	// If caller wants to keep old title/desc, they should pass empty string;
	// But better approach is for caller to indicate which flags were set and only update those.
	if task.Title == "" {
		task.Title = old.Title
	}
	if task.Description == "" {
		task.Description = old.Description
	}
	// For Done boolean, the caller must explicitly indicate whether to change it
	// (so this function will just set what caller provided)
	task.ID = old.ID
	s.Tasks[idx] = task

	data, err := json.MarshalIndent(s.Tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling tasks: %w", err)
	}
	if err := os.WriteFile(s.FilePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

// delete removes task with given ID
func (s *Store) delete(ID int) error {
	idx, err := s.GetTask(ID)
	if err != nil {
		return err
	}

	// remove element at idx
	s.Tasks = append(s.Tasks[:idx], s.Tasks[idx+1:]...)

	data, err := json.MarshalIndent(s.Tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling tasks: %w", err)
	}
	if err := os.WriteFile(s.FilePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func (s *Store) list() {
	fmt.Printf("ID\tTitle\tDescription\tCompleted\n")
	fmt.Printf("--\t-----\t-----------\t---------\n")
	for _, task := range s.Tasks {
		fmt.Printf("%d\t%s\t%s\t%t\n", task.ID, task.Title, task.Description, task.Done)
	}
}

func (s *Store) search(ID int) (*Task, error) {
	idx, err := s.GetTask(ID)
	if err != nil {
		return nil, err
	}
	return &s.Tasks[idx], nil
}

func loadStore(path string) (*Store, error) {
	s := &Store{FilePath: path, Tasks: []Task{}, nextID: 1}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&s.Tasks); err != nil {
		return nil, err
	}

	// compute nextID
	max := 0
	for _, task := range s.Tasks {
		if task.ID > max {
			max = task.ID
		}
	}
	s.nextID = max + 1

	return s, nil
}
