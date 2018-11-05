package main

import (
	"encoding/csv"
	"log"
	"os"
)

// Row as written to the CSV
type Row struct {
	Category    string
	Name        string
	Title       string
	Description string
	Notes       string
}

func (s *Service) newWriter(filename string) {
	s.wg.Add(1)
	defer s.wg.Done()

	file, err := os.Create(filename + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	s.writer = csv.NewWriter(file)
	defer s.writer.Flush()

	err = s.writer.Write([]string{
		"Category", "Name", "Title", "Description", "Notes"})
	if err != nil {
		log.Println(err)
		return
	}

	s.save = make(chan Row)

	for {
		row, more := <-s.save
		if !more {
			// channel closed, stop writer
			return
		}

		err = s.writer.Write([]string{
			row.Category,
			row.Name,
			row.Title,
			row.Description,
			row.Notes,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}