package main

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Creates a basic table with the headers provided
// and returns its table.Writer
// TODO: (#15:open) implement go-pretty tables for output
func createTable(headers ...string) table.Writer {
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	row := make(table.Row, len(headers))
	for i, h := range headers {
		row[i] = h
	}
	tw.AppendHeader(row)
	tw.SetStyle(table.StyleRounded)
	return tw
}
