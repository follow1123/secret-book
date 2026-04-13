package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BookManager struct {
	book *Book
}

func New(bookPath string) (*BookManager, error) {
	data, err := os.ReadFile(bookPath)
	if err != nil {
		return nil, fmt.Errorf("read book path: %s error:\n\t%w", bookPath, err)
	}
	book := &Book{}
	if err := json.Unmarshal(data, book); err != nil {
		return nil, fmt.Errorf("unmarshal book path: %s error:\n\t%w", bookPath, err)
	}

	return &BookManager{book: book}, nil
}
