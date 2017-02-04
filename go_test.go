package main

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestUtf8(t *testing.T) {
	fmt.Println(utf8.DecodeRuneInString("a世ssss"))
	fmt.Println(utf8.DecodeRuneInString("a"))
}
