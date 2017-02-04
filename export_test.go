package main

import (
	"testing"
)

func TestExport(testing *testing.T) {
	CheckInit()
	arguments := ParseArgs([]string{"--export", "--filter", "q>=2"})
	storage := NewFileStorage()
	HandleExport(storage, arguments)
}

func TestExport2(testing *testing.T) {
	CheckInit()
	arguments := ParseArgs([]string{"--export", "--filter", "q>=2", "--exclude-zh"})
	storage := NewFileStorage()
	HandleExport(storage, arguments)
}

func TestExportList(t *testing.T) {
	CheckInit()
	arguments := ParseArgs([]string{"--list", "--filter", "q>=1", "--exclude-zh"})
	storage := NewFileStorage()
	HandleList(storage, arguments)
}
