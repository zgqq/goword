package main

import (
	"testing"
)

func TestExport(testing *testing.T) {
	Init()
	arguments := ParseArgs([]string{"--export", "--filter", "q>=2"})
	storage := NewFileStorage()
	HandleExport(storage, arguments)
}

func TestExport2(testing *testing.T) {
	Init()
	arguments := ParseArgs([]string{"--export", "--filter", "q>=2", "--exclude-zh"})
	storage := NewFileStorage()
	HandleExport(storage, arguments)
}
