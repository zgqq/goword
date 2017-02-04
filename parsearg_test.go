package main

import (
	"testing"
)

func TestType(t *testing.T) {
	// // var storage Storage
	// // storage = &FileStorage{}
	// var storage interface{}
	// storage = 1
	// switch v := storage.(type) {
	// case FileStorage:

	// }
}

func TestGetArg(t *testing.T) {
	i := 0
	val, ok := getArgValue(&i, []string{"--export", "/tmp"})
	if !ok {
		t.Error("expect ok")
		t.FailNow()
	}
	if val != "/tmp" {
		t.Error("expect /tmp")
		t.FailNow()
	}
}

func TestParseExport(t *testing.T) {

	argument := ParseArgs([]string{"--export"})
	if argument.IsExport != true {
		t.Error("expect isExport to true")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--export", "/tmp"})
	if argument.ExportLoc != "/tmp" {
		t.Error("expect location /tmp")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--exclude-zh"})
	if argument.ExcludeZh != true {
		t.Error("expect exclude-zh true")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--filter", "q>2"})
	if argument.FilterCondi != nil && argument.FilterCondi.QCValue != 2 && argument.FilterCondi.IsQueryCount && argument.FilterCondi.Condi != GreaterThan {
		t.Error("expect great than 2")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--filter", "q>=3"})
	if argument.FilterCondi != nil && argument.FilterCondi.QCValue != 3 && argument.FilterCondi.IsQueryCount && argument.FilterCondi.Condi != GreaterOrEqual {
		t.Error("expect great or equal 3")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--filter", "q<2"})
	if argument.FilterCondi != nil && argument.FilterCondi.QCValue != 2 && argument.FilterCondi.IsQueryCount && argument.FilterCondi.Condi != LessThan {
		t.Error("expect less than 2")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--filter", "q<=3"})
	if argument.FilterCondi != nil && argument.FilterCondi.QCValue != 3 && argument.FilterCondi.IsQueryCount && argument.FilterCondi.Condi != LessOrEqual {
		t.Error("expect less or equal 3")
		t.FailNow()
	}

	argument = ParseArgs([]string{"--export", "/tmp", "--filter", "q<=3"})
	if argument.ExportLoc != "/tmp" {
		t.Error("expect location /tmp")
		t.FailNow()
	}
	if argument.FilterCondi != nil && argument.FilterCondi.QCValue != 3 && argument.FilterCondi.IsQueryCount && argument.FilterCondi.Condi != LessOrEqual {
		t.Error("expect less or equal 3")
		t.FailNow()
	}
}
