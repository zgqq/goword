package main

import (
	"testing"
)

func TestExport(testing *testing.T) {
	Init()
	ExportToTxtFile(NewFileStorage())
}
