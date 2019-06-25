package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/olekukonko/tablewriter"
	"io"
)

//go:generate mockery -name "ModuleGetter"

const ListModulesCommandString = "list modules"

type ModuleGetter interface {
	GetAll(ctx context.Context) (ch360.ModuleList, error)
}

type ListModules struct {
	client ModuleGetter
	writer io.Writer
}

func NewListModules(client ModuleGetter, out io.Writer) *ListModules {
	return &ListModules{
		client: client,
		writer: out,
	}
}

func (cmd *ListModules) Execute(ctx context.Context) error {
	modules, err := cmd.client.GetAll(ctx)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if len(modules) == 0 {
		fmt.Fprintln(cmd.writer, "No modules found.")
		return nil
	}

	table := tablewriter.NewWriter(cmd.writer)
	table.SetHeader([]string{"Name", "ID", "Summary"})
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetAutoWrapText(false)
	table.SetColumnSeparator("")

	for _, module := range modules {
		table.Append([]string{module.Name, module.ID, module.Summary})
	}
	table.Render()

	return nil
}

func (cmd ListModules) Usage() string {
	return ListModulesCommandString
}
