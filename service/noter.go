package service

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"
)

type Noter struct {
	basePath        string
	historyBasePath string
	saveHistory     bool
}

func NewNoter(basePath string, historyPath string, saveHistory bool) *Noter {
	if err := ensureDir(basePath); err != nil {
		panic(err)
	}
	if saveHistory {
		if err := ensureDir(historyPath); err != nil {
			panic(err)
		}
	}
	return &Noter{basePath: basePath, historyBasePath: historyPath, saveHistory: saveHistory}
}

func (s *Noter) getSavePath(noteId string) string {
	return path.Join(s.basePath, noteId+".json")
}

func (s *Noter) getHisSavePath(noteId string) string {
	return path.Join(s.historyBasePath, noteId+".json")
}

func (s *Noter) SaveNote(note *Note) error {
	note.ModifyTime = time.Now()

	err := saveJSON(note, s.getSavePath(note.Id))
	if err != nil {
		return fmt.Errorf("SaveNote: %v", err)
	}
	if s.saveHistory {
		if err := appendJSON(note, s.getHisSavePath(note.Id)); err != nil {
			return fmt.Errorf("SaveNote saveHistory: %v", err)
		}
	}
	return nil
}

func (s *Noter) GetNote(noteId string) (*Note, error) {
	path := s.getSavePath(noteId)
	note := Note{}
	err := readJSON(&note, path)
	if err != nil {
		return nil, fmt.Errorf("GetNote: %v", err)
	}
	return &note, nil
}

func (s *Noter) GetNoteHistory(noteId string) ([]Note, error) {
	lines, err := readLines(s.getHisSavePath(noteId))
	if err != nil {
		return nil, fmt.Errorf("GetNoteHistory: %v", err)
	}
	e := []error{}
	r := []Note{}
	for i, line := range lines {
		note := Note{}
		err := json.Unmarshal([]byte(line), &note)
		if err != nil {
			e = append(e, fmt.Errorf("line %d - %v", i, err))
		} else {
			r = append(r, note)
		}
	}
	var errs error = nil
	if len(e) > 0 {
		errs = e[0]
		for i := len(e) + 1; i < len(e); i++ {
			err = fmt.Errorf("%v; %v", err, e[i])
		}
	}
	return r, errs
}

func NewNote() *Note {
	note := Note{}
	note.Id = uuid.New().String()
	note.CreateTime = time.Now()
	return &note
}
