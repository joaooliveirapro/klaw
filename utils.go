package main

import "fmt"

// ANSI colors
const (
	Reset        = "\033[0m"
	Bold         = "\033[1m"
	Dim          = "\033[2m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Magenta      = "\033[35m"
	Cyan         = "\033[36m"
	Gray         = "\033[90m"
	GreenBGWhite = "\033[42;97m"
)

// color text using ANSI escape codes
func C(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// print error formatted
func IfErrPrint(err error) {
	if err != nil {
		fmt.Printf("%s %s", C("[error]", Red), err)
	}
}
