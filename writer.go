package main

import (
	"encoding/csv"
	"io"

	"github.com/olekukonko/tablewriter"
)

type writer interface {
	Write(...string) error
	Flush()
	Error() error
}

// ----------- CSV ----------- //
type csvWriter struct {
	w *csv.Writer
}

func NewCSVWriter(w io.Writer) csvWriter {
	return csvWriter{
		w: csv.NewWriter(w),
	}
}

func (cw csvWriter) Write(args ...string) error {
	return cw.w.Write(args)
}

func (cw csvWriter) Flush() {
	cw.w.Flush()
}

func (cw csvWriter) Error() error {
	return cw.w.Error()
}

// ----------- Table ----------- //
type tableWriter struct {
	w *tablewriter.Table
}

func NewTableWriter(w io.Writer, header []string) tableWriter {
	table := tablewriter.NewWriter(w)

	// Do not auto capitalize the header.
	table.SetAutoFormatHeaders(false)
	table.SetHeader(header)

	// If just one column, allow it to be extra long.
	if len(header) == 1 {
		table.SetAutoWrapText(false)
		table.SetReflowDuringAutoWrap(false)
	}

	return tableWriter{
		w: table,
	}
}

func (tw tableWriter) Write(args ...string) error {
	tw.w.Append(args)
	return nil
}

func (tw tableWriter) Flush() {
	tw.w.Render()
}

func (tw tableWriter) Error() error {
	return nil
}
