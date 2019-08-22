package commands

import (
	"github.com/olekukonko/tablewriter"
	"io"
)

// NewTable is a convenience method for creating tables used in a few places within the package.
func NewTable(writer io.Writer, headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetAutoWrapText(false)
	table.SetColumnSeparator("")

	return table
}
