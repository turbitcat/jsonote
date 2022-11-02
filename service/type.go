package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Note struct {
	Id         string
	Help       string
	CreateTime time.Time
	ModifyTime time.Time
	Content    any
	WriteKey   string
	ReadKey    string
}

type NoteHistory struct {
	Id    string
	Notes []Note
	Time  time.Time
}

func appendJSON(data any, path string) error {
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("appendJSON: %v", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("appendJSON: %v", err)
	}
	defer f.Close()
	j = append(j, '\n')
	if _, err = f.Write(j); err != nil {
		return fmt.Errorf("appendJSON: %v", err)
	}
	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func saveJSON(data any, path string) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0600)
	return err
}

func readJSON(data any, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("readJSON: %v", err)
	}
	err = json.Unmarshal(file, data)
	return err
}
