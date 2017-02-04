package main

import (
	"testing"
)

func TestInitStat(t *testing.T) {
	InitStat()
}

func TestLoadTodayStat(t *testing.T) {
	today := LoadTodayStat()
	if today == nil {
		t.Fail()
	}
	t.Log(today)
}

func TestIncrementTodayQueryCount(t *testing.T) {
	today := LoadTodayStat()
	addAfter := IncrementTodayQueryCount()
	if (today.QueryCount + 1) != addAfter.QueryCount {
		t.Log("expect querycount", (today.QueryCount + 1))
		t.Error("result querycount", today.QueryCount)
	}
	t.Log(today, addAfter)
}
